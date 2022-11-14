package lab

import (
	"dns-testbed-go/dns"
	"path/filepath"
	"strings"
	"time"
)

type timingExperiment struct {
	name     string
	zonesDir string
}

func newTimingExperiment(name, zonesDir string) *timingExperiment {
	return &timingExperiment{name: name, zonesDir: zonesDir}
}

func (t *timingExperiment) String() string {
	return t.name
}

func (t *timingExperiment) getMeasure() measure {
	return func(system *dns.System, delayInMS int) float64 {
		system.Target.SetZone(filepath.Join(t.zonesDir, t.name, "target.zone"))
		system.Target.SetDelay(time.Duration(delayInMS) * time.Millisecond)
		system.Client.Query("a1.target.com", "A", system.Resolver)
		targetLog := system.Target.ReadQueryLog(2 * time.Second)
		queryDuration, err := t.computeQueryDuration(targetLog)
		if err != nil {
			panic(err)
		}
		return queryDuration.Seconds()
	}
}

func (t *timingExperiment) computeQueryDuration(queryLog []byte) (time.Duration, error) {
	lines := strings.Split(string(queryLog), "\n")
	if len(lines) < 2 {
		return 0, nil
	}
	startTime, err := t.parseTimestamp(lines[0])
	if err != nil {
		return 0, err
	}
	endTime, err := t.parseTimestamp(lines[len(lines)-2])
	if err != nil {
		return 0, err
	}
	return endTime.Sub(startTime), nil
}

func (t *timingExperiment) parseTimestamp(queryLogLine string) (time.Time, error) {
	timestamp := strings.Split(queryLogLine, " ")[0] + " " + strings.Split(queryLogLine, " ")[1]
	return time.Parse("02-Jan-2006 15:04:05.000", timestamp)
}
