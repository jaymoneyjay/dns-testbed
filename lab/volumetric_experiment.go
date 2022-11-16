package lab

import (
	"dns-testbed-go/dns"
	"fmt"
	"path/filepath"
	"strings"
)

type VolumetricExperiment struct {
	name     string
	zonesDir string
}

func newVolumetricExperiment(name, zonesDir string) *VolumetricExperiment {
	return &VolumetricExperiment{name: name, zonesDir: zonesDir}
}

func (v *VolumetricExperiment) String() string {
	return v.name
}

func (v *VolumetricExperiment) getMeasure() measure {
	return func(system *dns.System, numberOfDelegations int, entryZone string) float64 {
		system.Inter.SetZone(v.getZonePath(numberOfDelegations, system.Inter.ID()))
		system.Target.SetZone(v.getZonePath(numberOfDelegations, system.Target.ID()))
		system.Client.Query(entryZone, "A", system.Resolver)
		targetLog := system.Target.ReadQueryLog(0)
		numberOfQueries := v.countQueries(targetLog)
		return numberOfQueries
	}
}

func (v *VolumetricExperiment) getZonePath(numberOfDelegations int, nsID string) string {
	return filepath.Join(v.zonesDir, v.name, nsID, fmt.Sprintf("ns-del-%d.zone", numberOfDelegations))
}

func (v *VolumetricExperiment) countQueries(queryLog []byte) float64 {
	lines := strings.Split(string(queryLog), "\n")
	return float64(len(lines))
}
