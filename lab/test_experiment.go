package lab

import (
	"dns-testbed-go/dns"
	"fmt"
	"path/filepath"
	"strings"
)

type testExperiment struct {
	name      string
	zonesDir  string
	entryZone string
}

func newTestExperiment(name, entryZone, zonesDir string) *testExperiment {
	return &testExperiment{name: name, entryZone: entryZone, zonesDir: zonesDir}
}

func (v *testExperiment) String() string {
	return v.name
}

func (v *testExperiment) getMeasure() measure {
	return func(system *dns.System, numberOfDelegations int) float64 {
		system.Inter.SetZone(v.getZonePath(numberOfDelegations, system.Inter.ID()))
		system.Target.SetZone(v.getZonePath(numberOfDelegations, system.Target.ID()))
		system.Client.Query(v.entryZone, "A", system.Resolver)
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
