package component

import (
	"dns-testbed-go/testbed/docker"
)

type Container struct {
	containerID string
	dockerCli   *docker.Client
}

func newContainer(containerID string) *Container {
	return &Container{
		containerID: containerID,
		dockerCli:   docker.NewClient(),
	}
}

func (c *Container) Start(implementation Implementation) error {
	return c.startBind9()
}

func (c *Container) exec(cmd []string) (docker.ExecResult, error) {
	return c.dockerCli.Exec(c.containerID, cmd)
}

func (c *Container) startBind9() error {
	_, err := c.exec([]string{"service", "bind9", "start"})
	return err
}
