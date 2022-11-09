package component

import "fmt"

type Implementation int

const (
	Bind_9_11_3 Implementation = iota
	Unbound17
	Unbound16
	Unbound10
	PowerDNS47
	Bind_9_18_4
)

func (i Implementation) String() string {
	return [...]string{"bind-9.11.3", "unbound-1.17.0", "unbound-1.16.0", "unbound-1.10.0", "powerDNS-4.7.3", "bind-9.18.4"}[i]
}

func (i Implementation) IP() string {
	return fmt.Sprintf("172.20.0.%d", 10+i)
}
