package component

import "fmt"

type Implementation int

const (
	Bind9 Implementation = iota
	Unbound17
	Unbound16
	Unbound10
	PowerDNS47
	Knot
)

func (i Implementation) String() string {
	return [...]string{"bind9", "unbound-1.17.0", "unbound-1.16.0", "unbound-1.10.0", "powerDNS-4.7", "knot"}[i]
}

func (i Implementation) IP() string {
	return fmt.Sprintf("172.20.0.%d", 10+i)
}
