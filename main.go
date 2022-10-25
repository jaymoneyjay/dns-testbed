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
	testbed.Client.SetResolver(component.Unbound17)
	queryResult, err := testbed.Query("target.com.", "A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
	queryResult, err = testbed.Query("inter.net.", "A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
	queryResult, err = testbed.Query("del.inter.net.", "A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(queryResult)
}
