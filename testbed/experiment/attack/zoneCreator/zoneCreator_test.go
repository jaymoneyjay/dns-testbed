package zoneCreator

import (
	"log"
	"reflect"
	"testing"
)

func TestShouldCreateZoneFragmentForNSDelegations(t *testing.T) {
	zoneCreator := NewZoneCreator([]string{"target.com.", "inter.net."})
	zoneCreator.AppendNSDelegations([]string{"entry"}, 10, 0, 1)
	actual := zoneCreator.zoneFragments
	expected := zoneFragmentMap{
		"inter.net.": newZoneFragment([]string{}),
		"target.com.": newZoneFragment([]string{
			"entry.target.com. IN NS del1.entry.inter.net.",
			"entry.target.com. IN NS del2.entry.inter.net.",
			"entry.target.com. IN NS del3.entry.inter.net.",
			"entry.target.com. IN NS del4.entry.inter.net.",
			"entry.target.com. IN NS del5.entry.inter.net.",
			"entry.target.com. IN NS del6.entry.inter.net.",
			"entry.target.com. IN NS del7.entry.inter.net.",
			"entry.target.com. IN NS del8.entry.inter.net.",
			"entry.target.com. IN NS del9.entry.inter.net.",
			"entry.target.com. IN NS del10.entry.inter.net.",
		}),
	}

	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldCreateLeafNodesForNSDelegations(t *testing.T) {
	zoneCreator := NewZoneCreator([]string{"source.com.", "inter.net.", "target.org."})
	actual := zoneCreator.AppendNSDelegations([]string{"entry1", "entry2"}, 2, 0, 2)

	expected := []string{
		"del1.entry1",
		"del2.entry1",
		"del1.entry2",
		"del2.entry2",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldCreateZoneFragmentForCNAMEChains(t *testing.T) {
	zoneCreator := NewZoneCreator([]string{"target.com.", "inter.net.", "inter.org."})
	zoneCreator.AppendCNAMEChain([]string{"entry1", "entry2"}, 7, []int{0, 1, 2})
	actual := zoneCreator.zoneFragments

	expected := zoneFragmentMap{
		"target.com.": newZoneFragment([]string{
			"entry1.target.com. IN CNAME chain2.entry1.inter.net.",
			"chain4.entry1.target.com. IN CNAME chain5.entry1.inter.net.",
			"entry2.target.com. IN CNAME chain2.entry2.inter.net.",
			"chain4.entry2.target.com. IN CNAME chain5.entry2.inter.net.",
		}),
		"inter.net.": newZoneFragment([]string{
			"chain2.entry1.inter.net. IN CNAME chain3.entry1.inter.org.",
			"chain5.entry1.inter.net. IN CNAME chain6.entry1.inter.org.",
			"chain2.entry2.inter.net. IN CNAME chain3.entry2.inter.org.",
			"chain5.entry2.inter.net. IN CNAME chain6.entry2.inter.org.",
		}),
		"inter.org.": newZoneFragment([]string{
			"chain3.entry1.inter.org. IN CNAME chain4.entry1.target.com.",
			"chain6.entry1.inter.org. IN CNAME chain7.entry1.target.com.",
			"chain3.entry2.inter.org. IN CNAME chain4.entry2.target.com.",
			"chain6.entry2.inter.org. IN CNAME chain7.entry2.target.com.",
		}),
	}
	if !actual.Equal(expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldCreateLeafNodesForCNAMEChains(t *testing.T) {
	zoneCreator := NewZoneCreator([]string{"target.com.", "inter.org.", "inter.com."})
	actual := zoneCreator.AppendCNAMEChain([]string{"entry1", "entry2", "entry3"}, 5, []int{0, 1})

	expected := []string{
		"chain5.entry1",
		"chain5.entry2",
		"chain5.entry3",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldCreateZoneFragmentForNSDelegationsWithCNAMEChains(t *testing.T) {
	zoneCreator := NewZoneCreator([]string{"target.com.", "inter.net.", "inter.org."})
	leafNodes := zoneCreator.AppendNSDelegations([]string{"entry"}, 3, 1, 2)
	zoneCreator.AppendCNAMEChain(leafNodes, 5, []int{2, 1, 0})
	actual := zoneCreator.zoneFragments

	expected := zoneFragmentMap{
		"target.com.": newZoneFragment([]string{
			"chain3.del1.entry.target.com. IN CNAME chain4.del1.entry.inter.org.",
			"chain3.del2.entry.target.com. IN CNAME chain4.del2.entry.inter.org.",
			"chain3.del3.entry.target.com. IN CNAME chain4.del3.entry.inter.org.",
		}),
		"inter.net.": newZoneFragment([]string{
			"entry.inter.net. IN NS del1.entry.inter.org.",
			"entry.inter.net. IN NS del2.entry.inter.org.",
			"entry.inter.net. IN NS del3.entry.inter.org.",
			"chain2.del1.entry.inter.net. IN CNAME chain3.del1.entry.target.com.",
			"chain2.del2.entry.inter.net. IN CNAME chain3.del2.entry.target.com.",
			"chain2.del3.entry.inter.net. IN CNAME chain3.del3.entry.target.com.",
		}),
		"inter.org.": newZoneFragment([]string{
			"del1.entry.inter.org. IN CNAME chain2.del1.entry.inter.net.",
			"del2.entry.inter.org. IN CNAME chain2.del2.entry.inter.net.",
			"del3.entry.inter.org. IN CNAME chain2.del3.entry.inter.net.",
			"chain4.del1.entry.inter.org. IN CNAME chain5.del1.entry.inter.net.",
			"chain4.del2.entry.inter.org. IN CNAME chain5.del2.entry.inter.net.",
			"chain4.del3.entry.inter.org. IN CNAME chain5.del3.entry.inter.net.",
		}),
	}
	if actual.Equal(expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldCreateLeafNodesForNSDelegationsWithCNAMEChains(t *testing.T) {
	zoneCreator := NewZoneCreator([]string{"target.com.", "inter.net.", "inter.org."})
	leafNodes := zoneCreator.AppendNSDelegations([]string{"entry"}, 3, 1, 2)
	actual := zoneCreator.AppendCNAMEChain(leafNodes, 5, []int{2, 1, 0})

	expected := []string{
		"chain5.del1.entry",
		"chain5.del2.entry",
		"chain5.del3.entry",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}
