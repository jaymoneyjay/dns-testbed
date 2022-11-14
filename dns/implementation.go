package dns

type implementation interface {
	Kind() string
	Version() string
	FlushCacheExecution() execution
	RestartExecution() execution
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

type execution struct {
	command              []string
	responseVerification func(string)
}
