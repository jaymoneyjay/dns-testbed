package main

import (
	"dns-testbed-go/lab"
)

func main() {
	dnsTestlab := lab.New("results")

	dnsTestlab.Conduct(
		lab.Test_QMIN,
		lab.NewDataIterator(getImplementations(), lab.MakeRange(10, 10, 1)),
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
