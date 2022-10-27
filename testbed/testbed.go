package testbed

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/docker"
)

type Testbed struct {
	Nameservers map[string][]*component.Nameserver
	Client      *component.Client
	Resolver    map[component.Implementation]*component.Resolver
}

func NewTestbed() (*Testbed, error) {
	root, err := component.NewNameserver("root", ".", "testbed/docker/buildContext/nameserver/root")
	if err != nil {
		return nil, err
	}
	com, err := component.NewNameserver("com", "com.", "testbed/docker/buildContext/nameserver/com")
	if err != nil {
		return nil, err
	}
	net, err := component.NewNameserver("net", "net.", "testbed/docker/buildContext/nameserver/net")
	if err != nil {
		return nil, err
	}
	target, err := component.NewNameserver("target-com", "target.com.", "testbed/docker/buildContext/nameserver/target-com")
	if err != nil {
		return nil, err
	}
	inter, err := component.NewNameserver("inter-net", "inter.net.", "testbed/docker/buildContext/nameserver/inter-net")
	if err != nil {
		return nil, err
	}
	nameservers := map[string][]*component.Nameserver{
		"root": {root},
		"tld": {
			com,
			net,
		},
		"sld": {
			target,
			inter,
		},
	}
	client, err := component.NewClient("client")
	if err != nil {
		return nil, err
	}
	bind9, err := component.NewResolver("resolver-bind9", "testbed/docker/buildContext/resolver/bind9")
	unbound17, err := component.NewResolver("resolver-unbound-1.17.0", "testbed/docker/buildContext/resolver/unbound-1.17.0")
	unbound16, err := component.NewResolver("resolver-unbound-1.16.0", "testbed/docker/buildContext/resolver/unbound-1.16.0")
	unbound10, err := component.NewResolver("resolver-unbound-1.10.0", "testbed/docker/buildContext/resolver/unbound-1.10.0")
	resolvers := map[component.Implementation]*component.Resolver{
		component.Bind9:     bind9,
		component.Unbound17: unbound17,
		component.Unbound16: unbound16,
		component.Unbound10: unbound10,
	}
	if err != nil {
		return nil, err
	}
	return &Testbed{
		Nameservers: nameservers,
		Client:      client,
		Resolver:    resolvers,
	}, nil
}

func (t *Testbed) Query(zone, record string) (docker.ExecResult, error) {
	return t.Client.Query(zone, record)
}

func (t *Testbed) CleanLogs() error {
	for _, resolver := range t.Resolver {
		err := resolver.CleanLog()
		if err != nil {
			return err
		}
	}
	for _, nsList := range t.Nameservers {
		for _, nameserver := range nsList {
			err := nameserver.CleanLog()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
