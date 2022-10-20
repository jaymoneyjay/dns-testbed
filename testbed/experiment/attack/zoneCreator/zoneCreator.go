package zoneCreator

type ZoneCreator struct {
	zones         []string
	zoneFragments zoneFragmentMap
}

func NewZoneCreator(zones []string) *ZoneCreator {
	zoneFragments := make(zoneFragmentMap)
	for _, zone := range zones {
		zoneFragments[zone] = newZoneFragment([]string{})
	}
	return &ZoneCreator{zones: zones, zoneFragments: zoneFragments}
}

func (z *ZoneCreator) AppendCNAMEChain(entryNodes []string, chainLength int, zoneIndices []int) []string {
	var zoneSubset []string
	for _, index := range zoneIndices {
		zoneSubset = append(zoneSubset, z.zones[index])
	}
	chain := newChain(chainLength, zoneSubset)
	var leafNodes []string
	for _, node := range entryNodes {
		zoneFragments, leafNodesSubSet := chain.Create(node)
		leafNodes = append(leafNodes, leafNodesSubSet...)
		z.mergeZoneFragments(zoneFragments)
	}
	return leafNodes
}

func (z *ZoneCreator) AppendNSDelegations(entryNodes []string, numberOfDelegations, indexSourceZone, indexTargetZone int) []string {
	delegation := newDelegation(numberOfDelegations, z.zones[indexSourceZone], z.zones[indexTargetZone])
	var leafNodes []string
	for _, node := range entryNodes {
		zoneFragments, leafNodesSubSet := delegation.Create(node)
		leafNodes = append(leafNodes, leafNodesSubSet...)
		z.mergeZoneFragments(zoneFragments)
	}
	return leafNodes
}

func (z *ZoneCreator) GetZoneFragment(zone string) string {
	return z.zoneFragments[zone].String()
}

func (z *ZoneCreator) mergeZoneFragments(zoneFragments zoneFragmentMap) {
	for zone, fragments := range zoneFragments {
		z.zoneFragments[zone] = z.zoneFragments[zone].merge(fragments)
	}
}
