package main

import (
	"dns-testbed-go/experiment"
	"dns-testbed-go/utils"
	"log"
)

func main() {
	subqueryCNAMEExperiment := experiment.NewExperiment(experiment.SubqueryCNAMEQMIN)
	err := subqueryCNAMEExperiment.Run(utils.MakeRange(1, 10, 1))
	if err != nil {
		log.Fatalln(err)
	}
}
