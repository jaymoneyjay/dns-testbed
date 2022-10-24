package main

import (
	"dns-testbed-go/testbed"
	"dns-testbed-go/testbed/component"
	"fmt"
	"log"
)

func main() {
	testbed, err := testbed.NewTestbed()
	if err != nil {
		log.Fatal(err)
	}
	err = testbed.CleanLogs()
	if err != nil {
		log.Fatal(err)
	}
	for _, sld := range testbed.Nameservers["sld"] {
		err := sld.SetZoneFile("subquery-unchained.zone")
		if err != nil {
			log.Fatal(err)
		}
	}
	// TODO: Find mistake why template.zone is still loaded
	err = testbed.Start(component.Bind9)
	if err != nil {
		log.Fatal(err)
	}
	queryResult, err := testbed.Query("entry.target.com.")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
	queryResult, err = testbed.Query("inter.net.")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
	queryResult, err = testbed.Query("del.inter.net.")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
}
