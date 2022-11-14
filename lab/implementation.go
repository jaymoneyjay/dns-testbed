package lab

type implementation int

const (
	Bind11 implementation = iota
	Bind18
	Unbound17
	Unbound16
	Unbound10
	PowerDNS
)

func (i implementation) ID() string {
	return []string{"bind-9.11.3", "bind-9.18.4", "unbound-1.17.0", "unbound-1.16.0", "unbound-1.10.0", "powerDNS-4.7.3"}[i]
}
