package lab

import (
	"dns-testbed-go/dns"
	"fmt"
	"path/filepath"
	"strings"
)

type VolumetricExperiment struct {
	name           string
	entryZone      string
	zonesDir       string
	zonePrefix     string
	intermediateNS bool
}

func newVolumetricExperiment(name, entryZone, zonesDir, zonePrefix string, intermediateNS bool) *VolumetricExperiment {
	return &VolumetricExperiment{name: name, entryZone: entryZone, zonesDir: zonesDir, zonePrefix: zonePrefix, intermediateNS: intermediateNS}
}

func (v *VolumetricExperiment) String() string {
	return v.name
}

func (v *VolumetricExperiment) getMeasure() measure {
	return func(system *dns.System, numberOfDelegations int) float64 {
		if v.intermediateNS {
			system.Inter.SetZone(v.getZonePath(numberOfDelegations, system.Inter.ID()))
		}
		system.Target.SetZone(v.getZonePath(numberOfDelegations, system.Target.ID()))
		system.Client.Query(v.entryZone, "A", system.Resolver)
		targetLog := system.Target.ReadQueryLog(0)
		numberOfQueries := v.countQueries(targetLog)
		return numberOfQueries
	}
}

func (v *VolumetricExperiment) getZonePath(numberOfDelegations int, nsID string) string {
	return filepath.Join(v.zonesDir, v.name, nsID, fmt.Sprintf("%s-%d.zone", v.zonePrefix, numberOfDelegations))
}

func (v *VolumetricExperiment) countQueries(queryLog []byte) float64 {
	lines := strings.Split(string(queryLog), "\n")
	return float64(len(lines))
}
