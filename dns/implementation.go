package dns

import "time"

type implementation interface {
	Kind() string
	Version() string
	flushCache(containerID string)
	restart(containerID string)
	readQueryLog(containerID, containerType string, minTimeout time.Duration) []byte
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
