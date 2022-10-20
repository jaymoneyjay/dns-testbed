package main

import (
	"dns-testbed-go/testbed"
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack"
	"fmt"
	"log"
)

func main() {
	testbed := testbed.NewTestbed()
	_, err := attack.NewTemplateAttack().WriteZoneFilesAndReturnEntryZone(0, testbed.Nameservers["sld"])
	if err != nil {
		log.Fatal(err)
	}
	err = testbed.Start(component.Bind9)
	if err != nil {
		log.Fatal(err)
	}
	queryResult, err := testbed.Query("target.com	")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
}
