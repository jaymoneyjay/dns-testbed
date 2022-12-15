package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Testbed struct {
	Templates string
	Build     string
	Zones     []*Zone
	Resolvers []*Resolver
	Client    *Client
	QMin      bool
}

type Zone struct {
	ID              string
	IP              string
	QName           string
	ZoneFileHost    string
	ZoneFileTarget  string
	Dir             string
	DefaultZoneFile string
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
	StartCommands   []string
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

func (c *Config) LoadTestbedConfig() (*Testbed, error) {
	c.v.SetConfigFile(filepath.Join(c.path, "testbed.yaml"))
	if err := c.v.ReadInConfig(); err != nil {
		return nil, err
	}
	testbedConfig := &Testbed{}
	if err := c.v.Unmarshal(testbedConfig); err != nil {
		return nil, err
	}
	c.completeZoneConfig(testbedConfig)
	c.completeResolverConfig(testbedConfig)
	testbedConfig.Client.Dir = filepath.Join(testbedConfig.Build, testbedConfig.Client.ID)
	return testbedConfig, nil
}

func (c *Config) completeZoneConfig(testbedConfig *Testbed) {
	for _, zone := range testbedConfig.Zones {
		zone.ID = generateZoneID(zone.QName)
		zone.ZoneFileHost = filepath.Join(testbedConfig.Build, "zones", fmt.Sprintf("%s.zone", zone.ID))
		zone.ZoneFileTarget = filepath.Join("/zones", fmt.Sprintf("%s.zone", zone.ID))
		zone.Dir = filepath.Join(testbedConfig.Build, zone.ID)
	}
}

func (c *Config) completeResolverConfig(testbedConfig *Testbed) {
	for _, resolver := range testbedConfig.Resolvers {
		impl := newImplementation(resolver.Implementation)
		resolver.ID = generateResolverID(resolver.Implementation, resolver.Version)
		resolver.Dir = filepath.Join(testbedConfig.Build, resolver.ID)
		resolver.ConfigTarget = impl.configTarget()
		resolver.RootHintsTarget = impl.rootHintsTarget()
		resolver.LogsTarget = impl.logsTarget()
		resolver.StartCommands = impl.startCommands()
	}
}

func generateZoneID(qname string) string {
	if qname == "." {
		return "root"
	}
	f := func(c rune) bool {
		return c == '.'
	}
	split := strings.FieldsFunc(qname, f)
	return strings.Join(split, "-")
}

func generateResolverID(impl, version string) string {
	return fmt.Sprintf("resolver-%s-%s", impl, version)
}
