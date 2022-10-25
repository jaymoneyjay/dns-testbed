package component

import "fmt"

type Implementation int

const (
	Bind9 Implementation = iota
	Unbound17
	PowerDNS
)

func (i Implementation) String() string {
	return [...]string{"bind9", "unbound-1.17.0", "powerDNS"}[i]
}

func (i Implementation) IP() string {
	return fmt.Sprintf("172.20.0.%d", 10+i)
}
