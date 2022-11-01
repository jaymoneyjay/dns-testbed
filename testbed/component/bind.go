package component

import (
	"dns-testbed-go/testbed/docker"
)

type bind struct {
	container *Container
}

func newBind(containerID string) (*bind, error) {
	container, err := NewContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &bind{container: container}, nil
}

func (b *bind) FlushCache() (docker.ExecResult, error) {
	return b.container.Exec([]string{"service", "bind9", "restart"})
}
