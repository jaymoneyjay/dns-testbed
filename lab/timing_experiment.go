package lab

import (
	"dns-testbed-go/dns"
	"path/filepath"
	"strings"
	"time"
)

type TimingExperiment struct {
	name      string
	zonesDir  string
	entryZone string
}

func newTimingExperiment(name, entryZone, zonesDir string) *TimingExperiment {
	return &TimingExperiment{name: name, entryZone: entryZone, zonesDir: zonesDir}
}

func (t *TimingExperiment) String() string {
	return t.name
}

func (t *TimingExperiment) getMeasure() measure {
	return func(system *dns.System, delayInMS int) float64 {
		system.Target.SetZone(filepath.Join(t.zonesDir, t.name, "target.zone"))
		system.Target.SetDelay(time.Duration(delayInMS) * time.Millisecond)
		t.warmup(system, delayInMS)
		system.Client.Query(t.entryZone, "A", system.Resolver)
		targetLog := system.Target.ReadQueryLog(2 * time.Second)
		queryDuration, err := t.computeQueryDuration(targetLog)
		if err != nil {
			panic(err)
		}
		return queryDuration.Seconds()
	}
}

func (t *TimingExperiment) warmup(system *dns.System, delayMS int) {
	zones := []string{
		"target.com",
		"www.target.com",
		"www1.target.com",
	}
	for _, zone := range zones {
		for i := 0; i < 10; i++ {
			system.Client.Query(zone, "A", system.Resolver)
		}
	}
	time.Sleep(time.Millisecond * 2 * time.Duration(delayMS))
	system.FlushQueryLogs()
}

func (t *TimingExperiment) computeQueryDuration(queryLog []byte) (time.Duration, error) {
	lines := strings.Split(string(queryLog), "\n")
	if len(lines) < 2 {
		return 0, nil
	}
	startTime, err := t.parseTimestamp(lines[0])
	if err != nil {
		return 0, err
	}
	endTime, err := t.parseTimestamp(lines[len(lines)-1])
	if err != nil {
		return 0, err
	}
	return endTime.Sub(startTime), nil
}

func (t *TimingExperiment) parseTimestamp(queryLogLine string) (time.Time, error) {
	elems := strings.Split(queryLogLine, " ")[0:2]
	timestamp := strings.Join(elems, " ")
	parsedTimestamp, err := time.Parse("02-Jan-2006 15:04:05.000", timestamp)
	if err == nil {
		return parsedTimestamp, nil
	}
	e := strings.Split(queryLogLine, " ")[0:3]
	timestamp = strings.Join(e, " ")
	return time.Parse("Jan 02 15:04:05", timestamp)
}
