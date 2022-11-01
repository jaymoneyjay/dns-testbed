package component

import "fmt"

type Implementation int

const (
	Bind9 Implementation = iota
	Unbound17
	Unbound16
	Unbound10
	PowerDNS47
)

func (i Implementation) String() string {
	return [...]string{"bind-9.11.3", "unbound-1.17.0", "unbound-1.16.0", "unbound-1.10.0", "powerDNS-4.7.3"}[i]
}

func (i Implementation) IP() string {
	return fmt.Sprintf("172.20.0.%d", 10+i)
}
