package dns

import (
	"dns-testbed-go/docker"
	"fmt"
)

type Client struct {
	ID        string
	dockerCli *docker.Client
}

func newClient(ID string, client *docker.Client) *Client {
	return &Client{
		ID:        ID,
		dockerCli: client,
	}
}

func (c *Client) Query(zone, record string, resolver *Resolver) string {
	execResult, err := c.dockerCli.Exec(c.ID, []string{"dig", "+tries=1", fmt.Sprintf("@%s", resolver.IP()), zone, record})
	if err != nil {
		panic(err)
	}
	return execResult.StdOut
}
