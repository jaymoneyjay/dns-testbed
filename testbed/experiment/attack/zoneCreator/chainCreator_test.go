package zoneCreator

import (
	"log"
	"reflect"
	"testing"
)

func TestShouldCreateChainLabelSequence(t *testing.T) {
	chain := newChain(10, []string{"target.com.", "inter1.com.", "inter2.com."})
	actual, _ := chain.createLabelSequence("ancestor")

	expected := []string{
		"ancestor.target.com.",
		"chain2.ancestor.inter1.com.",
		"chain3.ancestor.inter2.com.",
		"chain4.ancestor.target.com.",
		"chain5.ancestor.inter1.com.",
		"chain6.ancestor.inter2.com.",
		"chain7.ancestor.target.com.",
		"chain8.ancestor.inter1.com.",
		"chain9.ancestor.inter2.com.",
		"chain10.ancestor.target.com.",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

//TODO: test empty ancestor

func TestShouldCreateZoneFragmentsForSingleCNAMEChain(t *testing.T) {
	chain := newChain(10, []string{"target.com.", "inter1.com.", "inter2.com.", "inter3.com."})
	actual, _ := chain.Create("ancestor")

	expected := zoneFragmentMap{
		"target.com.": newZoneFragment([]string{
			"ancestor.target.com. IN CNAME chain2.ancestor.inter1.com.",
			"chain5.ancestor.target.com. IN CNAME chain6.ancestor.inter1.com.",
			"chain9.ancestor.target.com. IN CNAME chain10.ancestor.inter1.com.",
		}),
		"inter1.com.": newZoneFragment([]string{
			"chain2.ancestor.inter1.com. IN CNAME chain3.ancestor.inter2.com.",
			"chain6.ancestor.inter1.com. IN CNAME chain7.ancestor.inter2.com.",
		}),
		"inter2.com.": newZoneFragment([]string{
			"chain3.ancestor.inter2.com. IN CNAME chain4.ancestor.inter3.com.",
			"chain7.ancestor.inter2.com. IN CNAME chain8.ancestor.inter3.com.",
		}),
		"inter3.com.": newZoneFragment([]string{
			"chain4.ancestor.inter3.com. IN CNAME chain5.ancestor.target.com.",
			"chain8.ancestor.inter3.com. IN CNAME chain9.ancestor.target.com.",
		}),
	}

	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldCreateLeafNodesForSingleCNAMEChain(t *testing.T) {
	chain := newChain(3, []string{"target.com.", "inter.org."})
	_, actual := chain.Create("entry")

	expected := []string{"chain3.entry"}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}

}
