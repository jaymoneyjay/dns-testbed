package testbed

import (
	"errors"
	"fmt"
	"testbed/config"
)

type Resolver struct {
	Implementation
	*Container
	Dir string
}

func newResolver(resolverConfig *config.Resolver, templates string) *Resolver {
	container := NewContainer(resolverConfig.ID, resolverConfig.Dir, resolverConfig.IP)
	var implementation Implementation
	switch resolverConfig.Implementation {
	case "bind":
		implementation = newBind(templates, container)
	case "unbound":
		implementation = newUnbound(templates, container)
	case "powerdns":
		implementation = newPowerDNS(templates, container)
	default:
		panic(errors.New(fmt.Sprintf("Implementation kind %s has no instantiation method yet.", resolverConfig.Implementation)))
	}
	return &Resolver{
		Implementation: implementation,
		Container:      container,
	}
}
