package experiment

import (
	"dns-testbed-go/testbed"
	"dns-testbed-go/testbed/component"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"log"
	"os"
	"path/filepath"
	"time"
)

type TimingExperiment struct {
	attack     SlowDNS
	dnsTestbed *testbed.Testbed
	resultCSV  string
}

func NewTimingExperiment(attack SlowDNS) *TimingExperiment {
	t, err := testbed.NewTestbed()
	if err != nil {
		log.Fatal(err)
	}
	return &TimingExperiment{
		dnsTestbed: t,
		attack:     attack,
		resultCSV:  fmt.Sprintf("results/%s.csv", attack.String()),
	}
}

func (e *TimingExperiment) Run(delays []int, implementations []component.Implementation) error {
	err := e.dnsTestbed.CleanLogs()
	if err != nil {
		return err
	}
	target := e.dnsTestbed.Nameservers["sld"][0]
	var dataDelay []int
	var dataQueryDuration []float64
	var dataImpl []string
	zone := filepath.Join(e.attack.String(), "target.zone")
	err = target.SetZoneFile(zone)
	if err != nil {
		return err
	}
	for i, delay := range delays {
		err = e.dnsTestbed.FlushResolverCache()
		if err != nil {
			return err
		}
		delayResponse, err := target.SetDelay(delay)
		fmt.Println(delayResponse)
		fmt.Printf("Set delay at target to %d ms\n", delay)
		if err != nil {
			return err
		}
		for _, implementation := range implementations {
			e.dnsTestbed.Client.SetResolver(implementation)
			e.warmup("target.com", 3)
			e.warmup("www.target.com", 3)
			e.warmup("www1.target.com", 3)
			err = target.Clean()
			if err != nil {
				return err
			}
			_, err = e.dnsTestbed.Query("a1.target.com", "A")
			if err != nil {
				return err
			}
			queryDuration, err := target.GetQueryDuration(time.Duration(i) * 5 * time.Second)
			if err != nil {
				return err
			}
			fileName := fmt.Sprintf("%s-%dms.log", implementation.String(), delay)
			err = saveQueryLog(target, fileName, e.attack.String())
			if err != nil {
				return err
			}
			fmt.Printf("%s: %s\n", implementation.String(), queryDuration)
			if err != nil {
				return err
			}
			dataDelay = append(dataDelay, delay)
			dataQueryDuration = append(dataQueryDuration, queryDuration.Seconds())
			dataImpl = append(dataImpl, implementation.String())
		}
	}
	dfResults := dataframe.New(
		series.New(dataDelay, series.Int, "Delay"),
		series.New(dataQueryDuration, series.Int, "Query Duration"),
		series.New(dataImpl, series.String, "Implementation"),
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

func (e *TimingExperiment) warmup(zone string, repetitions int) {
	for i := 0; i < repetitions; i++ {
		_, _ = e.dnsTestbed.Query(zone, "A")
		time.Sleep(time.Second * 1)
	}
}
