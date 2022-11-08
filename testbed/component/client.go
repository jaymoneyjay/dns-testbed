package component

import (
	"dns-testbed-go/testbed/docker"
	"fmt"
)

type Client struct {
	*Container
	ResolverImplementation Implementation
}

func AttachClient(containerID string) (*Client, error) {
	container, err := NewContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &Client{
		Container:              container,
		ResolverImplementation: Bind9,
	}, nil
}

func (c *Client) SetResolver(resolver Implementation) {
	c.ResolverImplementation = resolver
}

func (c *Client) Query(zone, record string) (docker.ExecResult, error) {
	exec, err := c.Exec([]string{"dig", fmt.Sprintf("@%s", c.ResolverImplementation.IP()), zone, record})
	if err != nil {
		return docker.ExecResult{}, err
	}
	return exec, nil
}
