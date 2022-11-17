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
	runSubqueryExperiment(dnsTestLab, lab.Subquery_CNAME)
	runSubqueryExperiment(dnsTestLab, lab.Subquery_DNAME)
	runSubqueryExperiment(dnsTestLab, lab.Subquery_CNAME_Scrubbing)
	runSubqueryExperiment(dnsTestLab, lab.Subquery_DNAME_Scrubbing)
	runSlowDNSExperiment(dnsTestLab, lab.SlowDNS_CNAME_Scrubbing)
	runSlowDNSExperiment(dnsTestLab, lab.SlowDNS_DNAME_Scrubbing)
}

func runSubqueryExperiment(dnsTestLab *lab.Lab, experiment *lab.VolumetricExperiment) {
	dnsTestLab.Conduct(
		experiment,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
	)
	dnsTestLab.SaveResults()
}

func runSlowDNSExperiment(dnsTestLab *lab.Lab, experiment *lab.TimingExperiment) {
	dnsTestLab.Conduct(
		experiment,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(0, 1400, 200)),
	)
	dnsTestLab.SaveResults()
}
