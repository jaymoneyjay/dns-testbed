package component

import (
	"dns-testbed-go/testbed/docker"
	"fmt"
	"log"
	"os"
)

type Container struct {
	containerID string
	dockerCli   *docker.Client
	logger      *log.Logger
}

func newContainer(containerID string) (*Container, error) {
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
		containerID: containerID,
		dockerCli:   client,
		logger:      logger,
	}, nil
}

func (c *Container) exec(cmd []string) (docker.ExecResult, error) {
	execResult, err := c.dockerCli.Exec(c.containerID, cmd)
	c.logger.Println(execResult)
	return execResult, err
}

func (c *Container) startBind9() error {
	_, err := c.exec([]string{"service", "bind9", "start"})
	return err
}

func (c *Container) stopBind9() error {
	_, err := c.exec([]string{"service", "bind9", "stop"})
	return err
}

func (c *Container) restartBind9() error {
	_, err := c.exec([]string{"service", "bind9", "restart"})
	return err
}
