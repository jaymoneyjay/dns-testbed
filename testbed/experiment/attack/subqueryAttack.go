package attack

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack/zoneCreator"
	"fmt"
)

type Subquery struct {
	chainLength int
	entryNode   string
}

func NewSubquery(chainLength int) *Subquery {
	return &Subquery{chainLength: chainLength, entryNode: "entry"}
}

func (sq *Subquery) WriteZoneFilesAndReturnEntryZone(numberOfDelegations int, sldServers []*component.Nameserver) (string, error) {
	var zones []string
	for _, server := range sldServers {
		zones = append(zones, server.GetZone())
	}
	creator := zoneCreator.NewZoneCreator(zones)
	delegations := creator.AppendNSDelegations([]string{sq.entryNode}, numberOfDelegations, 0, 1)
	creator.AppendCNAMEChain(delegations, sq.chainLength, []int{1, 0})
	for _, ns := range sldServers {
		err := ns.WriteZone(creator.GetZoneFragment(ns.GetZone()), sq.createID(numberOfDelegations))
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s.%s", sq.entryNode, sldServers[0].GetZone()), nil
}

func (sq *Subquery) Name() string {
	return "subquery-unchained"
}

func (sq *Subquery) createID(numberOfDelegations int) string {
	return fmt.Sprintf("%s-chain%d-del%d.zone", sq.Name(), sq.chainLength, numberOfDelegations)
}
