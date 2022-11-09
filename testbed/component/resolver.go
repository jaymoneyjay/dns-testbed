package component

import (
	"dns-testbed-go/testbed/docker"
	"errors"
	"fmt"
)

type Resolver interface {
	FlushCache() (docker.ExecResult, error)
}

type Version int

func AttachResolver(implementation Implementation) (Resolver, error) {
	containerID := fmt.Sprintf("resolver-%s", implementation.String())
	switch implementation {
	case Bind_9_11_3:
		return newBind(containerID)
	case Unbound10, Unbound16, Unbound17:
		return newUnbound(containerID)
	case PowerDNS47:
		return newPowerDNS(containerID)
	default:
		return nil, errors.New(fmt.Sprintf("no instantiation method for implementation %s", implementation.String()))
	}
}
