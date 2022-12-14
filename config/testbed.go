package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	v *viper.Viper
}

func New() *Config {
	v := viper.New()
	v.AddConfigPath("config")
	v.SetDefault("Build", "build")
	v.SetDefault("Templates", filepath.Join("testbed", "templates"))
	v.SetDefault("QMin", false)
	return &Config{
		v: v,
	}
}

func (c *Config) SetConfig(config string) error {
	configSrc, err := os.ReadFile(config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("config", "config.yaml"), configSrc, 0777); err != nil {
		return err
	}
	return nil
}
func (c *Config) LoadTestbedConfig() (*Testbed, error) {
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

func (c *Config) Load(key string) interface{} {
	return c.v.Get(key)
}
