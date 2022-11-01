package component

import (
	"dns-testbed-go/testbed/docker"
)

type unbound struct {
	container *Container
}

func newUnbound(containerID string) (*unbound, error) {
	container, err := NewContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &unbound{
		container: container,
	}, nil
}

func (u *unbound) FlushCache() (docker.ExecResult, error) {
	return u.container.Exec([]string{"unbound-control", "reload"})
}
