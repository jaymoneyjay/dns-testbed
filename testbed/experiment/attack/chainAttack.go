package attack

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack/zoneCreator"
	"fmt"
)

type chainAttack struct {
	entryNode   string
	chainLength int
}

func NewChainAttack(chainLength int) *chainAttack {
	return &chainAttack{entryNode: "entry", chainLength: chainLength}
}

func (c *chainAttack) WriteZoneFilesAndReturnEntryZone(param int, sldServers []*component.Nameserver) (string, error) {
	var zones []string
	for _, server := range sldServers {
		zones = append(zones, server.GetZone())
	}
	creator := zoneCreator.NewZoneCreator(zones)
	creator.AppendCNAMEChain([]string{c.entryNode}, c.chainLength, []int{0, 1})
	for _, ns := range sldServers {
		err := ns.WriteZone(creator.GetZoneFragment(ns.GetZone()), c.Name())
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s.%s", c.entryNode, sldServers[0].GetZone()), nil
}

func (c *chainAttack) Name() string {
	return "chain"
}
