package main

import (
	"path/filepath"
	"testbed/config"
	"testbed/testbed"
	"testing"
)

var testbedConfig = &config.Testbed{
	Zones: []*config.Zone{
		{
			QName: ".",
			IP:    "172.20.0.2",
		},
		{
			QName: "com.",
			IP:    "172.20.0.3",
		},
	},
	Resolvers: []*config.Resolver{
		{
			Implementation: "bind",
			Version:        "9.18.4",
			IP:             "172.20.0.51",
		},
	},
	Client: &config.Client{
		ID: "client",
		IP: "172.20.0.99",
	},
}

var zones = filepath.Join("zones", "default")

func TestShouldStartNewTestbed(t *testing.T) {
	testbed.New(testbedConfig).Start()
}

func TestShouldStopTestbed(t *testing.T) {
	testbed.New(testbedConfig).Stop()
}
