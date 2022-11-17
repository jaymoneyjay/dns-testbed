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
	Subquery_CNAME                = newVolumetricExperiment("subquery+CNAME", "del.inter.net.", "zones")
	Subquery_CNAME_Scrubbing      = newVolumetricExperiment("subquery+CNAME+scrubbing", "del.inter.net.", "zones")
	Subquery_CNAME_Scrubbing_QMIN = newVolumetricExperiment("subquery+CNAME+scrubbing+QMIN", "del.inter.net.", "zones")
	Subquery_DNAME                = newVolumetricExperiment("subquery+DNAME", "del.inter.net.", "zones")
	Subquery_DNAME_Scrubbing      = newVolumetricExperiment("subquery+DNAME+scrubbing", "del.inter.net.", "zones")
	SlowDNS_CNAME                 = newTimingExperiment("slowDNS+CNAME", "a1.target.com.", "zones")
	SlowDNS_CNAME_Scrubbing       = newTimingExperiment("slowDNS+CNAME+scrubbing", "a1.target.com.", "zones")
	SlowDNS_DNAME_Scrubbing       = newTimingExperiment("slowDNS+DNAME+scrubbing", "sub.a1.target.com.", "zones")
	SlowDNS_CNAME_Scrubbing_QMIN  = newTimingExperiment("slowDNS+CNAME+scrubbing", "a1.target.com.", "zones")
	Test_QMIN                     = newTestExperiment("test+QMIN", "a1.target.com.", "zones")
)
