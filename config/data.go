package config

type Testbed struct {
	Templates string
	Build     string
	Zones     []*Zone
	Resolvers []*Resolver
	Client    *Client
	QMin      bool
}

type Zone struct {
	ID             string
	IP             string
	QName          string
	ZoneFileHost   string
	ZoneFileTarget string
	Dir            string
}

type Resolver struct {
	ID              string
	IP              string
	Version         string
	Implementation  string
	Dir             string
	ConfigTarget    string
	RootHintsTarget string
	LogsTarget      string
}

type Client struct {
	ID         string
	IP         string
	Nameserver string
	Dir        string
}

type DockerCompose struct {
	Zones     []*Zone
	Resolvers []*Resolver
	Client    *Client
}

type ResolverConfig struct {
	QMin string
}
