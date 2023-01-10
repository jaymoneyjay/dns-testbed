package config

import (
	"path/filepath"
	"time"
)

type Experiment struct {
	Name         string
	ResolverIDs  []string
	ZonesDir     string
	Delay        []int
	DelayedZones []string
	Target       string
	Measure      string
	Query        *query
	Warmup       []string
	QMin         bool
	Dest         string
	SaveLogs     bool
	Timeout      time.Duration
}

type query struct {
	RecordType string
	QName      string
}

func (c *Config) LoadExperimentConfig(path string) (*Experiment, error) {
	c.v.SetDefault("saveLogs", true)
	c.v.SetDefault("delay", []int{0})
	c.v.SetDefault("dest", filepath.Join("validation", "results"))
	c.v.SetConfigFile(path)
	if err := c.v.ReadInConfig(); err != nil {
		return nil, err
	}
	experimentConfig := &Experiment{}
	if err := c.v.Unmarshal(experimentConfig); err != nil {
		return nil, err
	}
	return experimentConfig, nil
}
