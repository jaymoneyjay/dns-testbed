package main

import (
	"dns-testbed-go/testbed"
	"dns-testbed-go/testbed/component"
	"fmt"
	"log"
)

func main() {
	t, err := testbed.NewTestbed()
	if err != nil {
		log.Fatal(err)
	}
	err = t.CleanLogs()
	if err != nil {
		log.Fatal(err)
	}
	target := t.Nameservers["sld"][0]
	inter := t.Nameservers["sld"][1]
	err = target.SetZoneFile("subquery-unchained-20.zone")
	if err != nil {
		log.Fatal(err)
	}
	t.Client.SetResolver(component.Unbound10)
	var results []int
	for i := 11; i < 21; i++ {
		err = inter.SetZoneFile(fmt.Sprintf("subquery-unchained-%d.zone", i))
		if err != nil {
			log.Fatal(err)
		}
		_, err := t.Client.Query("target.com", "A")
		if err != nil {
			log.Fatal(err)
		}
		_, err = t.Client.Query("inter.net", "A")
		if err != nil {
			log.Fatal(err)
		}
		err = target.CleanLog()
		if err != nil {
			log.Fatal(err)
		}
		//time.Sleep(5 * time.Second)
		queryResult, err := t.Client.Query("del.inter.net", "A")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(queryResult)
		numberOfQueries, err := target.CountQueries()
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, numberOfQueries)
		fmt.Println(numberOfQueries)
	}
	fmt.Print(results)
}
