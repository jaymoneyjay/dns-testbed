package experiment

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"strings"
	"testbed/config"
	"testbed/testbed"
	"time"
)

type Experiment struct {
	config *config.Experiment
}

func New(experimentConfig *config.Experiment) *Experiment {
	return &Experiment{config: experimentConfig}
}

func (e *Experiment) Run(testbed *testbed.Testbed) error {
	testbed.Reset()
	var dataResolver []string
	var dataZone []string
	var dataDelay []int
	var dataResult []int
	perm := os.FileMode(0777)
	resultDir := filepath.Join(e.config.Dest, e.config.Name)
	if _, err := os.Stat(resultDir); err == nil {
		if err := os.RemoveAll(resultDir); err != nil {
			return err
		}
	}
	if err := os.Mkdir(resultDir, perm); err != nil {
		return err
	}
	logDir := filepath.Join(resultDir, "logs")
	if err := os.Mkdir(logDir, perm); err != nil {
		return err
	}
	isVolume := e.config.Measure == "volume"
	isDuration := e.config.Measure == "duration"
	resolverIDs := e.config.ResolverIDs
	if resolverIDs == nil {
		for _, resolver := range testbed.Resolvers {
			resolverIDs = append(resolverIDs, resolver.ID)
		}
	}
	zoneConfigurations, err := os.ReadDir(e.config.ZonesDir)
	if err != nil {
		return err
	}
	bar := progressbar.Default(
		int64(
			len(resolverIDs)*
				len(zoneConfigurations)*
				len(e.config.Delay),
		),
		fmt.Sprintf("run %s experiment", e.config.Name),
	)
	for _, resolverID := range resolverIDs {
		resolver, err := testbed.FindResolver(resolverID)
		if err != nil {
			return err
		}
		resolver.SetConfig(e.config.QMin, true)
		for _, zoneConfig := range zoneConfigurations {
			if strings.HasPrefix(zoneConfig.Name(), ".") {
				err := bar.Add(len(e.config.Delay))
				if err != nil {
					return err
				}
				continue
			}
			testbed.SetZoneFiles(filepath.Join(e.config.ZonesDir, zoneConfig.Name()))
			for _, delay := range e.config.Delay {
				testbed.SetDelay(time.Duration(delay)*time.Millisecond, e.config.DelayedZones)
				testbed.Flush()
				if isDuration {
					e.warmup(testbed, resolverID)
				}
				testbed.FlushQueryLogs()
				testbed.Query(resolverID, e.config.Query.QName, e.config.Query.RecordType)
				result, _ := testbed.Measure(isVolume, isDuration, e.config.Target)
				dataResolver = append(dataResolver, resolverID)
				dataZone = append(dataZone, zoneConfig.Name())
				dataDelay = append(dataDelay, delay)
				dataResult = append(dataResult, int(result))
				if e.config.SaveLogs {
					currentLogDir := filepath.Join(
						logDir,
						fmt.Sprintf(
							"r-%s-z-%s-d-%d-qmin-%t",
							resolverID,
							zoneConfig.Name(),
							delay,
							e.config.QMin,
						),
					)
					if err := os.Mkdir(currentLogDir, perm); err != nil {
						return err
					}
					testbed.SaveLogs(resolverID, currentLogDir)
				}
				err := bar.Add(1)
				if err != nil {
					return err
				}
			}
		}
	}
	dfResult := dataframe.New(
		series.New(dataResolver, series.String, "resolver"),
		series.New(dataZone, series.String, "zone"),
		series.New(dataDelay, series.Int, "delay"),
		series.New(dataResult, series.Int, "result"),
	)
	resultsFile, err := os.Create(
		filepath.Join(resultDir, "data.csv"),
	)
	if err != nil {
		return err
	}
	if err := dfResult.WriteCSV(resultsFile); err != nil {
		return err
	}
	return nil
}

func (e *Experiment) warmup(testbed *testbed.Testbed, resolverID string) {
	for _, qname := range e.config.Warmup {
		for i := 0; i < 3; i++ {
			testbed.Query(resolverID, qname, "A")
		}
	}
	time.Sleep(1 * time.Second)
	testbed.FlushQueryLogs()
}
