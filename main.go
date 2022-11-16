package main

import (
	"dns-testbed-go/lab"
)

func main() {
	dnsTestLab := lab.New("results")
	runSlowDNSExperiment(dnsTestLab, lab.SlowDNS_CNAME_Scrubbing)
}

func runSubqueryExperiment(dnsTestLab *lab.Lab, experiment *lab.VolumetricExperiment) {
	dnsTestLab.Conduct(
		experiment,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
		"del.inter.net",
	)
	dnsTestLab.SaveResults()
}

func runSlowDNSExperiment(dnsTestLab *lab.Lab, experiment *lab.TimingExperiment) {
	dnsTestLab.Conduct(
		experiment,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(800, 800, 200)),
		"a1.target.com",
	)
	dnsTestLab.SaveResults()
}

func getImplementations() []string {
	return []string{
		//"bind-9.18.4",
		//"unbound-1.10.0",
		//"unbound-1.16.0",
		"powerDNS-4.7.3",
	}
}
