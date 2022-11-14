package dns

import "dns-testbed-go/docker"

type System struct {
	Target    *Nameserver
	Inter     *Nameserver
	Resolver  *Resolver
	Client    *Client
	resolvers map[string]*Resolver
}

func New() *System {
	dockerCli := docker.NewClient()
	resolvers := map[string]*Resolver{
		"bind-9.11.3":    newResolver(Bind, "9.11.3", "172.20.0.10", dockerCli),
		"unbound-1.17.0": newResolver(Unbound, "1.17.0", "172.20.0.11", dockerCli),
		"unbound-1.16.0": newResolver(Unbound, "1.16.0", "172.20.0.12", dockerCli),
		"unbound-1.10.0": newResolver(Unbound, "1.10.0", "172.20.0.13", dockerCli),
		"powerDNS-4.7.3": newResolver(PowerDNS, "4.7.3", "172.20.0.14", dockerCli),
		"bind-9.18.4":    newResolver(Bind, "9.18.4", "172.20.0.15", dockerCli),
	}
	return &System{
		Target:    newNameserver("target-com", dockerCli),
		Inter:     newNameserver("inter-net", dockerCli),
		resolvers: resolvers,
		Client:    newClient("client", dockerCli),
		Resolver:  resolvers["bind-9.11.3"],
	}
}

func (s *System) SetResolver(resolverID string) {
	s.Resolver = s.resolvers[resolverID]
	s.Resolver.FlushCache()
}

func (s *System) FlushQueryLogs() {
	s.Target.flushQueryLog()
	s.Inter.flushQueryLog()
	s.Resolver.flushQueryLog()
}
