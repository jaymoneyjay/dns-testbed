package main

import (
	"dns-testbed-go/testbed"
	"dns-testbed-go/testbed/component"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"log"
	"os"
	"time"
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

	var dataNSDel []int
	var dataNumQueries []int
	var dataImpl []string
	for implementation := range t.Resolver {
		impl, nsDel, numQueries := validate(t, implementation)
		dataNSDel = append(dataNSDel, nsDel...)
		dataNumQueries = append(dataNumQueries, numQueries...)
		dataImpl = append(dataImpl, impl...)
	}
	dfResults := dataframe.New(
		series.New(dataImpl, series.String, "Implementation"),
		series.New(dataNumQueries, series.Int, "Amplification"),
		series.New(dataNSDel, series.Int, "NS Delegations"),
	)
	resultsFile, err := os.Create("results.csv")
	if err != nil {
		log.Fatal(err)
	}
	err = dfResults.WriteCSV(resultsFile)
	if err != nil {
		log.Fatal(err)
	}
}

func validate(t *testbed.Testbed, implementation component.Implementation) ([]string, []int, []int) {
	target := t.Nameservers["sld"][0]
	inter := t.Nameservers["sld"][1]
	err := target.SetZoneFile("subquery-unchained-20.zone")
	if err != nil {
		log.Fatal(err)
	}
	t.Client.SetResolver(implementation)
	_, err = t.Client.Query("target.com", "A")
	if err != nil {
		log.Fatal(err)
	}
	_, err = t.Client.Query("inter.net", "A")
	if err != nil {
		log.Fatal(err)
	}
	var dataNSDel []int
	var dataNumQueries []int
	var dataImpl []string
	for i := 4; i < 11; i++ {
		err = inter.SetZoneFile(fmt.Sprintf("subquery-unchained-%d.zone", i))
		if err != nil {
			log.Fatal(err)
		}
		err = target.CleanLog()
		if err != nil {
			log.Fatal(err)
		}
		queryResult, err := t.Client.Query("del.inter.net", "A")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(queryResult)
		numberOfQueries, err := target.CountQueries()
		fmt.Println(numberOfQueries)
		if err != nil {
			log.Fatal(err)
		}
		dataNSDel = append(dataNSDel, i)
		dataNumQueries = append(dataNumQueries, numberOfQueries)
		dataImpl = append(dataImpl, implementation.String())
		time.Sleep(1 * time.Second)
	}
	return dataImpl, dataNSDel, dataNumQueries
}
