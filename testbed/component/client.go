package component

import (
	"dns-testbed-go/testbed/docker"
	"fmt"
)

type Client struct {
	*Container
	resolver Implementation
}

func NewClient(containerID string) (*Client, error) {
	container, err := newContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &Client{
		Container: container,
		resolver:  Bind9,
	}, nil
}

func (c *Client) SetResolver(resolver Implementation) {
	c.resolver = resolver
}

func (c *Client) Query(zone, record string) (docker.ExecResult, error) {
	return c.exec([]string{"dig", fmt.Sprintf("@%s", c.resolver.IP()), zone, record})
}
