package main

import (
	"dns-testbed-go/lab"
)

func main() {
	dnsTestlab := lab.New("results")

	dnsTestlab.Conduct(
		lab.SlowDNS_CNAME_Scrubbing,
		lab.NewDataIterator([]string{"powerDNS-4.7.3"}, lab.MakeRange(0, 800, 200)),
	)
	dnsTestlab.SaveResults()
}

func getImplementations() []string {
	return []string{
		"bind-9.11.3",
		"unbound-1.10.0",
		"unbound-1.16.0",
		"powerDNS-4.7.3",
	}
}
