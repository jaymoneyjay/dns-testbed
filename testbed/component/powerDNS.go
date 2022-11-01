package component

import (
	"dns-testbed-go/testbed/docker"
)

type powerDNS struct {
	container *Container
}

func newPowerDNS(containerID string) (*powerDNS, error) {
	container, err := NewContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &powerDNS{container: container}, nil
}

func (p *powerDNS) FlushCache() (docker.ExecResult, error) {
	return p.container.Exec([]string{"/etc/init.d/pdns-recursor", "restart"})
}
