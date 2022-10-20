package main

import (
	"dns-testbed-go/testbed"
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack"
	"log"
)

func main() {
	testbed := testbed.NewTestbed()
	zone, _ := attack.NewTemplateAttack().WriteZoneFilesAndReturnEntryZone(0, testbed.Nameservers["sld"])
	err := testbed.Start(component.Bind9)
	if err != nil {
		log.Fatal(err)
	}
	err = testbed.Query(zone)
	if err != nil {
		log.Fatal(err)
	}
}
