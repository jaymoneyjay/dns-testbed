package lab

import (
	"dns-testbed-go/dns"
)

type measure func(system *dns.System, x int) float64

type experiment interface {
	getMeasure() measure
	String() string
}

var (
	Subquery_CNAME                = newVolumetricExperiment("subquery+CNAME", "zones")
	Subquery_CNAME_Scrubbing      = newVolumetricExperiment("subquery+CNAME+scrubbing", "zones")
	Subquery_CNAME_Scrubbing_QMIN = newVolumetricExperiment("subquery+CNAME+scrubbing+QMIN", "zones")
	Subquery_DNAME                = newVolumetricExperiment("subquery+DNAME", "zones")
	Subquery_DNAME_Scrubbing      = newVolumetricExperiment("subquery+DNAME+scrubbing", "zones")
	SlowDNS_CNAME                 = newTimingExperiment("slowDNS+CNAME", "zones")
	SlowDNS_CNAME_Scrubbing       = newTimingExperiment("slowDNS+CNAME+Scrubbing", "zones")
	SlowDNS_CNAME_Scrubbing_QMIN  = newTimingExperiment("slowDNS+CNAME+Scrubbing", "zones")
	Test_QMIN                     = newTestExperiment("test+QMIN", "zones")
)
