package zoneCreator

import (
	"dns-testbed-go/testbed/experiment/attack/labelCreator"
	"fmt"
)

type delegation struct {
	numberOfDelegations int
	sourceZone          string
	targetZone          string
	dnsKeyword          string
}

func newDelegation(numberOfDelegations int, sourceZone, targetZone string) *delegation {
	return &delegation{numberOfDelegations: numberOfDelegations, dnsKeyword: "NS", sourceZone: sourceZone, targetZone: targetZone}
}

// Create generates delegations according to the delegation parameters and returns the RRs in a zoneFragmentMap together with a list of the current leaf nodes
func (d *delegation) Create(entryPoint string) (zoneFragmentMap, []string) {
	labelSequence, leafNodes := d.createLabelSequence(entryPoint)
	records := d.createRecords(labelSequence)
	return d.createZoneFragments(records), leafNodes
}

func (d *delegation) createLabelSequence(entryPoint string) ([]string, []string) {
	delegationLabel := labelCreator.NewDelegationLabel()
	var labels []string
	var leafNodes []string
	sourceLabel := fmt.Sprintf("%s.%s", entryPoint, d.sourceZone)
	labels = append(labels, sourceLabel)
	for i := 0; i < d.numberOfDelegations; i++ {
		targetLabel := delegationLabel.Step()
		targetLabel = fmt.Sprintf("%s.%s.%s", targetLabel, entryPoint, d.targetZone)
		labels = append(labels, targetLabel)
		leafNodes = append(leafNodes, fmt.Sprintf("%s.%s", delegationLabel, entryPoint))
	}
	return labels, leafNodes
}

func (d *delegation) createRecords(labels []string) []string {
	var records []string
	for i := 1; i < len(labels); i++ {
		records = append(records, fmt.Sprintf("%s IN %s %s", labels[0], d.dnsKeyword, labels[i]))
	}
	return records
}

func (d *delegation) createZoneFragments(records []string) zoneFragmentMap {
	fragments := zoneFragmentMap{
		d.sourceZone: newZoneFragment(records),
	}
	return fragments
}
