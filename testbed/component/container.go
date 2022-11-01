package component

import (
	"dns-testbed-go/testbed/docker"
	"fmt"
	"log"
	"os"
)

type Container struct {
	ContainerID string
	DockerCli   *docker.Client
	logger      *log.Logger
}

func NewContainer(containerID string) (*Container, error) {
	client, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	logFile, err := os.Create(fmt.Sprintf("testbed/component/logs/%s.log", containerID))
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile, "DOCKER: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &Container{
		ContainerID: containerID,
		DockerCli:   client,
		logger:      logger,
	}, nil
}

func (c *Container) Exec(cmd []string) (docker.ExecResult, error) {
	execResult, err := c.DockerCli.Exec(c.ContainerID, cmd)
	c.logger.Println(execResult)
	return execResult, err
}

func (c *Container) startBind9() error {
	_, err := c.Exec([]string{"service", "bind9", "start"})
	return err
}

func (c *Container) stopBind9() error {
	_, err := c.Exec([]string{"service", "bind9", "stop"})
	return err
}

func (c *Container) restartBind9() error {
	_, err := c.Exec([]string{"service", "bind9", "restart"})
	return err
}

func (c *Container) Restart() error {
	return c.DockerCli.RestartContainer(c.ContainerID)
}
