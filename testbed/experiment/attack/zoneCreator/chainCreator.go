package zoneCreator

import (
	"dns-testbed-go/testbed/experiment/attack/labelCreator"
	"fmt"
)

type chain struct {
	chainLength int
	dnsKeyword  string
	zones       []string
}

func newChain(chainLength int, zones []string) *chain {
	return &chain{chainLength: chainLength, dnsKeyword: "CNAME", zones: zones}
}

// Create generates dns chains according to the chain parameters and returns the RRs in a zoneFragmentMap together with a list of the current leaf nodes
func (c *chain) Create(entryPoint string) (zoneFragmentMap, []string) {
	labelSequence, leafNodes := c.createLabelSequence(entryPoint)
	records := c.createRecords(labelSequence)
	return c.splitIntoZoneFragments(records), leafNodes
}

func (c *chain) createLabelSequence(entryPoint string) ([]string, []string) {
	chainLabel := labelCreator.NewChainLabel()
	labels := []string{
		fmt.Sprintf("%s.%s", entryPoint, c.zones[0]),
	}
	chainLabel.Step()
	for i := 1; i < c.chainLength; i++ {
		label := chainLabel.Step()
		label = fmt.Sprintf("%s.%s.%s", label, entryPoint, c.zones[i%len(c.zones)])
		labels = append(labels, label)
	}
	leafNodes := []string{fmt.Sprintf("%s.%s", chainLabel, entryPoint)}
	return labels, leafNodes
}

func (c *chain) createRecords(labels []string) []string {
	var records []string
	for i := 0; i < len(labels)-1; i++ {
		records = append(records, fmt.Sprintf("%s IN %s %s", labels[i], c.dnsKeyword, labels[i+1]))
	}
	return records
}

func (c *chain) splitIntoZoneFragments(records []string) zoneFragmentMap {
	var fragments = make(zoneFragmentMap)
	for _, zone := range c.zones {
		fragments[zone] = newZoneFragment([]string{})
	}
	for i, record := range records {
		fragments[c.zones[i%len(c.zones)]].records = append(fragments[c.zones[i%len(c.zones)]].records, record)
	}
	return fragments
}
