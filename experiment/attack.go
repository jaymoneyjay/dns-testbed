package experiment

type SubqueryUnchained int

const (
	SubqueryBasic SubqueryUnchained = iota
	SubqueryCNAME
	SubqueryCNAME_QMIN
)

func (i SubqueryUnchained) String() string {
	return [...]string{"subquery", "subquery+CNAME", "subquery+CNAME+QMIN"}[i]
}

type SlowDNS int

const (
	SlowDNS_Basic SlowDNS = iota
	SlowDNS_CNAME
	SlowDNS_CNAME_QMIN
)

func (i SlowDNS) String() string {
	return [...]string{"slowDNS", "slowDNS+CNAME", "slowDNS+CNAME+QMIN"}[i]
}
