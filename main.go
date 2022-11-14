package main

import (
	"dns-testbed-go/lab"
)

func main() {
	dnsTestlab := lab.New("results")

	dnsTestlab.Conduct(
		lab.Subquery_CNAME,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
	)
	dnsTestlab.SaveResults()

	dnsTestlab.Conduct(
		lab.Subquery_CNAME_Scrubbing,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
	)
	dnsTestlab.SaveResults()

	dnsTestlab.Conduct(
		lab.Subquery_CNAME_Scrubbing_QMIN,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
	)
	dnsTestlab.SaveResults()

	dnsTestlab.Conduct(
		lab.Subquery_DNAME,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
	)
	dnsTestlab.SaveResults()

	dnsTestlab.Conduct(
		lab.Subquery_DNAME_Scrubbing,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(1, 10, 1)),
	)
	dnsTestlab.SaveResults()

	/*dnsTestlab.Conduct(
		lab.SlowDNS_CNAME_Scrubbing,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(0, 1400, 200)),
	)
	dnsTestlab.SaveResults()*/
}

func getImplementations() []string {
	return []string{
		"bind-9.11.3",
		"unbound-1.10.0",
		"unbound-1.16.0",
		"powerDNS-4.7.3",
	}
}
