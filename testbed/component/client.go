package component

import "dns-testbed-go/testbed/docker"

type Client struct {
	*Container
}

func NewClient(containerID string) (*Client, error) {
	container, err := newContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &Client{
		Container: container,
	}, nil
}

func (c *Client) Query(zone string) (docker.ExecResult, error) {
	return c.exec([]string{"dig", zone})
}
