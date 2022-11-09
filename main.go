package main

import (
	"dns-testbed-go/experiment"
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/utils"
	"log"
)

func main() {
	runTimingExperiment(experiment.SlowDNS_CNAME)
	//runSubqueryExperiment(experiment.SubqueryCNAME)
}

func runSubqueryExperiment(attack experiment.SubqueryUnchained) {
	subqueryCNAMEExperiment := experiment.NewSubqueryExperiment(attack)
	err := subqueryCNAMEExperiment.Run(utils.MakeRange(1, 10, 1))
	if err != nil {
		log.Fatalln(err)
	}
}

func runTimingExperiment(attack experiment.SlowDNS) {
	implementations := []component.Implementation{
		component.Bind_9_11_3,
		//component.Unbound10,
		component.Unbound16,
		component.PowerDNS47,
	}
	timingExperiment := experiment.NewTimingExperiment(attack)
	err := timingExperiment.Run(utils.MakeRange(0, 1400, 200), implementations)
	if err != nil {
		log.Fatalln(err)
	}
}
