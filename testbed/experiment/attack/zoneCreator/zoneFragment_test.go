package zoneCreator

import (
	"log"
	"testing"
)

func TestShouldResultInEqual(t *testing.T) {
	if !(zoneFragmentMap{
		"key1": newZoneFragment([]string{"a1", "a2"}),
		"key2": newZoneFragment([]string{"b1", "b2"}),
		"key3": newZoneFragment([]string{"c1", "c2"}),
		"key4": newZoneFragment([]string{"c1", "d2"}),
	}.Equal(zoneFragmentMap{
		"key1": newZoneFragment([]string{"a1", "a2"}),
		"key2": newZoneFragment([]string{"b1", "b2"}),
		"key3": newZoneFragment([]string{"c1", "c2"}),
		"key4": newZoneFragment([]string{"c1", "d2"}),
	})) {
		log.Fatalf("Should be equal.")
	}
}

func TestShouldIgnoreOrderAndResultInEqual(t *testing.T) {
	if !(zoneFragmentMap{
		"key1": newZoneFragment([]string{"a1", "a2"}),
		"key2": newZoneFragment([]string{"b1", "b2"}),
		"key3": newZoneFragment([]string{"c1", "c2"}),
		"key4": newZoneFragment([]string{"c1", "d2"}),
	}.Equal(zoneFragmentMap{
		"key1": newZoneFragment([]string{"a1", "a2"}),
		"key2": newZoneFragment([]string{"b1", "b2"}),
		"key4": newZoneFragment([]string{"c1", "d2"}),
		"key3": newZoneFragment([]string{"c1", "c2"}),
	})) {
		log.Fatalf("Should be equal.")
	}
}

func TestShouldRecognizeNotEqual(t *testing.T) {
	if (zoneFragmentMap{
		"key1": newZoneFragment([]string{"a1", "a2"}),
		"key2": newZoneFragment([]string{"b1", "b2"}),
		"key3": newZoneFragment([]string{"c1", "c2"}),
		"key4": newZoneFragment([]string{"c1", "d2"}),
	}.Equal(zoneFragmentMap{
		"key1": newZoneFragment([]string{"a1", "a2"}),
		"key2": newZoneFragment([]string{"b1", "b2"}),
		"key3": newZoneFragment([]string{"c1", "c2"}),
	})) {
		log.Fatalf("Should not be equal.")
	}
}

func TestShouldReturnString(t *testing.T) {
	actual := newZoneFragment([]string{"a1", "a2"}).String()

	expected := "a1\na2\n"
	if actual != expected {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}
