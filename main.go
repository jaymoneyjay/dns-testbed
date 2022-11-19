package main

import (
	"dns-testbed-go/lab"
)

func getImplementations() []string {
	return []string{
		"bind-9.18.4",
		"unbound-1.10.0",
		"unbound-1.16.0",
		"powerDNS-4.7.3",
	}
}

func main() {
	dnsTestLab := lab.New("results")
	/*
		runVolumetricExperiment(dnsTestLab, lab.Subquery_CNAME, []string{"bind-9.18.4"})
		runVolumetricExperiment(dnsTestLab, lab.Subquery_DNAME, getImplementations())

			runVolumetricExperiment(dnsTestLab, lab.Subquery_CNAME_Scrubbing, getImplementations())
			runVolumetricExperiment(dnsTestLab, lab.Subquery_DNAME_Scrubbing, getImplementations())
			runTimingExperiment(dnsTestLab, lab.SlowDNS_CNAME_Scrubbing, getImplementations())
			runTimingExperiment(dnsTestLab, lab.SlowDNS_DNAME_Scrubbing, getImplementations())
	*/
	runVolumetricExperiment(dnsTestLab, lab.Chain_CNAME_Scrubbing, []string{"bind-9.18.4"})
}

func runVolumetricExperiment(dnsTestLab *lab.Lab, experiment *lab.VolumetricExperiment, implementations []string) {
	dnsTestLab.Conduct(
		experiment,
		lab.NewDataIterator(implementations, lab.MakeRange(9, 18, 1)),
	)
	dnsTestLab.SaveResults()
}

func runTimingExperiment(dnsTestLab *lab.Lab, experiment *lab.TimingExperiment, implementations []string) {
	dnsTestLab.Conduct(
		experiment,
		lab.NewDataIterator(implementations, lab.MakeRange(0, 1400, 200)),
	)
	dnsTestLab.SaveResults()
}
