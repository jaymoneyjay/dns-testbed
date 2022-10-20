package zoneCreator

import (
	"fmt"
	"reflect"
)

type ZoneFragment struct {
	records []string
}

func newZoneFragment(records []string) *ZoneFragment {
	return &ZoneFragment{records: records}
}

func (zf *ZoneFragment) merge(other *ZoneFragment) *ZoneFragment {
	return &ZoneFragment{records: append(zf.records, other.records...)}
}

func (zf *ZoneFragment) String() string {
	result := ""
	for _, record := range zf.records {
		result += fmt.Sprintln(record)
	}
	return result
}

type zoneFragmentMap map[string]*ZoneFragment

func (zfs zoneFragmentMap) String() string {
	result := ""
	for zone, fragment := range zfs {
		result += fmt.Sprintf("%s: %s\n", zone, fragment)
	}
	return result
}

func (zfs zoneFragmentMap) Equal(other zoneFragmentMap) bool {
	equal := true
	for zone, fragment := range zfs {
		if !reflect.DeepEqual(other[zone], fragment) {
			equal = false
		}
	}
	return equal
}
