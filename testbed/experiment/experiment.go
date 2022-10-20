package experiment

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack"
)

type Experiment struct {
	target          component.Logging
	attack          attack.Attack
	dataRange       []int
	implementations []component.Implementation
}

func NewExperiment(attack attack.Attack, dataRange []int) *Experiment {
	return &Experiment{
		attack:    attack,
		dataRange: dataRange,
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
		_, err = client.Query(entryZone)
		if err != nil {
			return nil, nil, err
		}
		numberOfQueries, err := targetComponent.CountQueries()
		if err != nil {
			return nil, nil, err
		}
		data = append(data, dataPoint)
		values = append(values, numberOfQueries)
	}
	return data, values, nil
}
