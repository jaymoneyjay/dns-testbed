package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"regexp"
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

func (c *Client) Query(zone, record string, resolver *Resolver) error {
	execResult, err := c.dockerCli.Exec(c.ID, []string{"dig", fmt.Sprintf("@%s", resolver.IP()), zone, record})
	if err != nil {
		panic(err)
	}
	matched, err := regexp.MatchString(";; QUESTION SECTION:", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !matched {
		return errors.New(fmt.Sprintf("could not query zone %s for %s record:\n%s", zone, record, execResult.StdOut))
	}
	return nil
}
