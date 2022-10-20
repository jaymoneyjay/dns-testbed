package zoneCreator

import (
	"log"
	"reflect"
	"testing"
)

// TODO: test for empty entry point
func TestShouldGenerateZoneFragmentForSingleNSDelegation(t *testing.T) {
	delegation := newDelegation(10, "source.com.", "target.com.")
	actual, _ := delegation.Create("entry")
	expected := zoneFragmentMap{
		"source.com.": newZoneFragment([]string{
			"entry.source.com. IN NS del1.entry.target.com.",
			"entry.source.com. IN NS del2.entry.target.com.",
			"entry.source.com. IN NS del3.entry.target.com.",
			"entry.source.com. IN NS del4.entry.target.com.",
			"entry.source.com. IN NS del5.entry.target.com.",
			"entry.source.com. IN NS del6.entry.target.com.",
			"entry.source.com. IN NS del7.entry.target.com.",
			"entry.source.com. IN NS del8.entry.target.com.",
			"entry.source.com. IN NS del9.entry.target.com.",
			"entry.source.com. IN NS del10.entry.target.com.",
		})}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}

func TestShouldGenerateLeafNodesForSingleNSDelegation(t *testing.T) {
	delegation := newDelegation(4, "source.com.", "target.org.")
	_, actual := delegation.Create("entry")

	expected := []string{
		"del1.entry",
		"del2.entry",
		"del3.entry",
		"del4.entry",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}
