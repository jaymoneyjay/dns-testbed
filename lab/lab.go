package lab

import (
	"dns-testbed-go/dns"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type lab struct {
	dnsSystem    *dns.System
	dataX        []int
	dataY        []float64
	dataHue      []string
	logs         []*queryLog
	resultsDir   string
	experimentID string
}

func New(resultsDir string) *lab {
	return &lab{dnsSystem: dns.New(), logs: []*queryLog{}, resultsDir: resultsDir}
}

func (l *lab) Conduct(experiment experiment, dataIterator *dataIterator) {
	l.experimentID = fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02-15:04:05"), experiment)
	fmt.Printf("# Start measurements for %s experiment:\n", experiment)
	l.reset()
	for dataIterator.hasNextDataPoint() {
		hue, x := dataIterator.getNextDataPoint()
		l.dnsSystem.SetResolver(hue)
		l.dnsSystem.FlushQueryLogs()
		y := experiment.getMeasure()(l.dnsSystem, x)
		l.appendResult(x, y, hue)
		l.appendLogs(x, hue)
		fmt.Printf("\t%s, %d, %f\n", hue, x, y)
	}
}

func (l *lab) reset() {
	l.logs = []*queryLog{}
	l.dataX = []int{}
	l.dataY = []float64{}
	l.dataHue = []string{}
	l.dnsSystem.Target.SetDelay(0)
	l.dnsSystem.Inter.SetDelay(0)
}

func (l *lab) appendResult(dataX int, dataY float64, dataHue string) {
	l.dataX = append(l.dataX, dataX)
	l.dataY = append(l.dataY, dataY)
	l.dataHue = append(l.dataHue, dataHue)
}

func (l *lab) appendLogs(x int, hue string) {
	targetLog := l.dnsSystem.Target.ReadQueryLog(0)
	interLog := l.dnsSystem.Inter.ReadQueryLog(0)
	resolverLog := l.dnsSystem.Resolver.ReadQueryLog(0)
	l.dnsSystem.FlushQueryLogs()
	queryLog := newQueryLog(
		targetLog,
		interLog,
		resolverLog,
		hue,
		x,
	)
	l.logs = append(l.logs, queryLog)
}

func (l *lab) SaveResults() {
	perm := fs.ModePerm
	createDirIfNotExists(l.resultsDir)
	experimentResultsDir := filepath.Join(l.resultsDir, l.experimentID)
	err := os.Mkdir(experimentResultsDir, perm)
	if err != nil {
		panic(err)
	}
	experimentLogsDir := filepath.Join(experimentResultsDir, "logs")
	err = os.Mkdir(experimentLogsDir, perm)
	if err != nil {
		panic(err)
	}
	l.saveData(experimentResultsDir)
	l.saveLogs(experimentLogsDir)
}

func (l *lab) saveData(dir string) {
	dfResults := dataframe.New(
		series.New(l.dataX, series.Int, "x"),
		series.New(l.dataY, series.Float, "y"),
		series.New(l.dataHue, series.String, "hue"),
	)
	resultsFile, err := os.Create(
		filepath.Join(dir, "data.csv"),
	)
	if err != nil {
		panic(err)
	}
	err = dfResults.WriteCSV(resultsFile)
	if err != nil {
		panic(err)
	}
}

func (l *lab) saveLogs(dir string) {
	for _, log := range l.logs {
		log.save(dir)
	}
}

func createDirIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, fs.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
