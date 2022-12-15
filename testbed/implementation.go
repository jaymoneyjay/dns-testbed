package testbed

type Implementation interface {
	flushCache()
	reload()
	start()
	filterQueries(queryLog []byte) []byte
	SetConfig(qmin, reload bool)
}

type implementationKind int

const (
	Bind implementationKind = iota
	Unbound
	PowerDNS
)

func (i implementationKind) String() string {
	return []string{
		"bind",
		"unbound",
		"powerDNS",
	}[i]
}
