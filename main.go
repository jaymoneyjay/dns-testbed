package main

import (
	"dns-testbed-go/lab"
)

func main() {
	dnsTestlab := lab.New("results")

	dnsTestlab.Conduct(
		lab.SlowDNS_CNAME_Scrubbing,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(0, 1400, 200)),
		"a1.target.com",
	)
	dnsTestlab.SaveResults()
}

func getImplementations() []string {
	return []string{
		"bind-9.18.4",
		//"unbound-1.10.0",
		"unbound-1.16.0",
		"powerDNS-4.7.3",
	}
}
