package experiment

import (
	"dns-testbed-go/testbed"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"log"
	"os"
	"path/filepath"
)

type SubqueryExperiment struct {
	attack     SubqueryUnchained
	dnsTestbed *testbed.Testbed
	resultCSV  string
}

func NewSubqueryExperiment(attack SubqueryUnchained) *SubqueryExperiment {
	t, err := testbed.NewTestbed()
	if err != nil {
		log.Fatal(err)
	}
	return &SubqueryExperiment{
		dnsTestbed: t,
		attack:     attack,
		resultCSV:  fmt.Sprintf("results/%s.csv", attack.String()),
	}
}

func (e *SubqueryExperiment) Run(nsDelegations []int) error {
	err := e.dnsTestbed.CleanLogs()
	if err != nil {
		return err
	}
	target := e.dnsTestbed.Nameservers["sld"][0]
	inter := e.dnsTestbed.Nameservers["sld"][1]
	_, err = target.SetDelay(0)
	if err != nil {
		return err
	}
	_, err = inter.SetDelay(0)
	if err != nil {
		return err
	}
	var dataNSDel []int
	var dataNumQueries []int
	var dataImpl []string
	for _, nsDel := range nsDelegations {
		err = e.dnsTestbed.FlushResolverCache()
		if err != nil {
			return err
		}
		for implementation := range e.dnsTestbed.Resolver {
			e.dnsTestbed.Client.SetResolver(implementation)
			zone := filepath.Join(e.attack.String(), fmt.Sprintf("ns-del-%d.zone", nsDel))
			err = inter.SetZoneFile(zone)
			if err != nil {
				return err
			}
			err = target.SetZoneFile(zone)
			if err != nil {
				return err
			}
			err = target.Clean()
			if err != nil {
				return err
			}
			queryResult, err := e.dnsTestbed.Query("del.inter.net", "A")
			if err != nil {
				return err
			}
			fmt.Print(queryResult)
			numberOfQueries, err := target.CountQueries()
			fmt.Println(numberOfQueries)
			if err != nil {
				return err
			}
			err = saveQueryLog(target, fmt.Sprintf("%s-%d.log", implementation.String(), nsDel), e.attack.String())
			if err != nil {
				return err
			}
			dataNSDel = append(dataNSDel, nsDel)
			dataNumQueries = append(dataNumQueries, numberOfQueries)
			dataImpl = append(dataImpl, implementation.String())
		}
	}
	dfResults := dataframe.New(
		series.New(dataImpl, series.String, "Implementation"),
		series.New(dataNumQueries, series.Int, "Amplification"),
		series.New(dataNSDel, series.Int, "NS Delegations"),
	)
	resultsFile, err := os.Create(e.resultCSV)
	if err != nil {
		return err
	}
	err = dfResults.WriteCSV(resultsFile)
	if err != nil {
		return err
	}
	return nil
}
