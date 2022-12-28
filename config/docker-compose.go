package config

type DockerCompose struct {
	Zones     []*Zone
	Resolvers []*Resolver
	Client    *Client
}
