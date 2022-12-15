package config

import "path/filepath"

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

var TestExperiment = &Experiment{
	Name:         "test-experiment",
	ResolverIDs:  []string{"resolver-bind-9.18.4", "resolver-unbound-1.10.0"},
	ZonesDir:     filepath.Join("zones", "CNAME+scrubbing"),
	Delay:        []int{0},
	DelayedZones: nil,
	Target:       "target-com",
	Measure:      "volume",
	Query: &query{
		RecordType: "",
		QName:      "a1.target.com",
	},
	Warmup:   nil,
	QMin:     false,
	Dest:     "results",
	SaveLogs: true,
}
