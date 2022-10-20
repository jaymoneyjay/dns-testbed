package experiment

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/experiment/attack"
	"dns-testbed-go/utils"
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestShouldCreateDataFromSubqueriesUnchainedExperiment(t *testing.T) {
	client := component.NewClient("client")
	sldServers := []*component.Nameserver{
		component.NewNameserver("target-com", "target.com.", "../docker/buildContext/target-com"),
		component.NewNameserver("inter-net", "inter.net.", "../docker/buildContext/inter-net"),
	}
	actual, _, err := NewExperiment(
		attack.NewSubquery(10),
		utils.MakeRange(1, 2, 1),
	).Run(client, *sldServers[0], sldServers)
	if err != nil {
		panic(err)
	}

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf(fmt.Sprint("Actual:", actual, "\n", "Expected:", expected))
	}
}

func TestShouldCreateValuesFromSubqueriesUnchainedExperiment(t *testing.T) {
	client := component.NewClient("client")
	sldServers := []*component.Nameserver{
		component.NewNameserver("target-com", "target.com.", "../docker/buildContext/target-com"),
		component.NewNameserver("inter-net", "inter.net.", "../docker/buildContext/inter-net"),
	}
	subqueryAttack := attack.NewSubquery(3)
	_, actual, err := NewExperiment(
		subqueryAttack,
		utils.MakeRange(1, 10, 1),
	).Run(client, *sldServers[0], sldServers)
	if err != nil {
		panic(err)
	}

	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf(fmt.Sprint("Actual:", actual, "\n", "Expected:", expected))
	}
}
