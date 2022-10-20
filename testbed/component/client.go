package component

type Client struct {
	*Container
}

func NewClient(containerID string) *Client {
	return &Client{
		Container: newContainer(containerID),
	}
}

func (c *Client) Query(zone string) error {
	return c.dockerCli.Exec(c.containerID, []string{"dig", zone})
}
