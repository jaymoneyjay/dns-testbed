package lab

import (
	"dns-testbed-go/dns"
	"fmt"
	"path/filepath"
	"strings"
)

type testExperiment struct {
	name     string
	zonesDir string
}

func newTestExperiment(name, zonesDir string) *testExperiment {
	return &testExperiment{name: name, zonesDir: zonesDir}
}

func (v *testExperiment) String() string {
	return v.name
}

func (v *testExperiment) getMeasure() measure {
	return func(system *dns.System, numberOfDelegations int) float64 {
		system.Inter.SetZone(v.getZonePath(numberOfDelegations, system.Inter.ID()))
		system.Target.SetZone(v.getZonePath(numberOfDelegations, system.Target.ID()))
		err := system.Client.Query("lab12.lab11.lab10.lab9.lab8.lab7.lab6.lab5.lab4.lab3.lab2.lab1.a1.target.com.", "A", system.Resolver)
		if err != nil {
			panic(err)
		}
		targetLog := system.Target.ReadQueryLog(0)
		numberOfQueries := v.countQueries(targetLog)
		return numberOfQueries
	}
}

func (v *testExperiment) getZonePath(numberOfDelegations int, nsID string) string {
	return filepath.Join(v.zonesDir, v.name, nsID, fmt.Sprintf("ns-del-%d.zone", numberOfDelegations))
}

func (v *testExperiment) countQueries(queryLog []byte) float64 {
	lines := strings.Split(string(queryLog), "\n")
	return float64(len(lines) - 1)
}
