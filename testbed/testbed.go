package testbed

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/docker"
)

type Testbed struct {
	Nameservers map[string][]*component.Nameserver
	Client      *component.Client
	Resolver    map[component.Implementation]component.Resolver
}

func NewTestbed() (*Testbed, error) {
	root, err := component.AttachRoot()
	if err != nil {
		return nil, err
	}
	com, err := component.AttachCOM()
	if err != nil {
		return nil, err
	}
	net, err := component.AttachNET()
	if err != nil {
		return nil, err
	}
	target, err := component.AttachTarget()
	if err != nil {
		return nil, err
	}
	inter, err := component.AttachInter()
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
	client, err := component.AttachClient("client")
	if err != nil {
		return nil, err
	}
	bind9, err := component.AttachResolver(component.Bind9)
	unbound16, err := component.AttachResolver(component.Unbound16)
	unbound10, err := component.AttachResolver(component.Unbound10)
	powerDNS47, err := component.AttachResolver(component.PowerDNS47)
	resolvers := map[component.Implementation]component.Resolver{
		component.Bind9:      bind9,
		component.Unbound16:  unbound16,
		component.Unbound10:  unbound10,
		component.PowerDNS47: powerDNS47,
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
	for _, nsList := range t.Nameservers {
		for _, nameserver := range nsList {
			err := nameserver.Clean()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Testbed) FlushResolverCache() error {
	for _, resolver := range t.Resolver {
		_, err := resolver.FlushCache()
		if err != nil {
			return err
		}
	}
	return nil
}
