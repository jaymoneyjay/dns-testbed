package experiment

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack"
	"log"
	"os"
)

type Experiment struct {
	attack      attack.Attack
	dataRange   []int
	queryLogger *log.Logger
}

func NewExperiment(attack attack.Attack, dataRange []int) *Experiment {
	logWriter, err := os.Open("logs/query.log")
	if err != nil {
		log.Fatal(err)
	}

	return &Experiment{
		attack:      attack,
		dataRange:   dataRange,
		queryLogger: log.New(logWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (e *Experiment) Run(client *component.Client, targetComponent component.Logging, sldServers []*component.Nameserver) ([]int, []int, error) {
	var data []int
	var values []int
	for _, dataPoint := range e.dataRange {
		err := targetComponent.CleanLog()
		if err != nil {
			return nil, nil, err
		}
		entryZone, err := e.attack.WriteZoneFilesAndReturnEntryZone(dataPoint, sldServers)
		if err != nil {
			return nil, nil, err
		}
		queryResponse, err := client.Query(entryZone)
		if err != nil {
			return nil, nil, err
		}
		e.queryLogger.Print(queryResponse)
		numberOfQueries, err := targetComponent.CountQueries()
		if err != nil {
			return nil, nil, err
		}
		data = append(data, dataPoint)
		values = append(values, numberOfQueries)
	}
	return data, values, nil
}
