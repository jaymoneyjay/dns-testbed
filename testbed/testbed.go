package testbed

import (
	component2 "dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/docker"
	"dns-testbed-go/testbed/experiment"
	"dns-testbed-go/utils"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type Testbed struct {
	Nameservers     map[string][]*component2.Nameserver
	Client          *component2.Client
	Resolver        *component2.Resolver
	dockerClient    *docker.Client
	implementations []component2.Implementation
}

func NewTestbed() *Testbed {
	return &Testbed{Nameservers: map[string][]*component2.Nameserver{
		"root": {component2.NewNameserver("root", ".", "docker/buildContext/root")},
		"tld": {
			component2.NewNameserver("com", "com.", "docker/buildContext/com"),
			component2.NewNameserver("net", "net.", "docker/buildContext/net"),
		},
		"sld": {
			component2.NewNameserver("target-com", "target.com.", "docker/buildContext/target-com"),
			component2.NewNameserver("inter-net", "inter.net.", "docker/buildContext/target-net"),
		},
	},
		Client:          component2.NewClient("client"),
		Resolver:        component2.NewResolver("resolver", "docker/buildContext/resolver"),
		dockerClient:    docker.NewClient(),
		implementations: []component2.Implementation{component2.Bind9},
	}
}

func (t *Testbed) Start(implementation component2.Implementation) error {
	for _, ns := range t.Nameservers {
		for _, c := range ns {
			err := c.Start(implementation)
			if err != nil {
				return err
			}
		}
	}
	return t.Resolver.Start(implementation)
}

func (t *Testbed) Query(zone string) error {
	return t.Client.Query(zone)
}

func (t *Testbed) Run(experiment *experiment.Experiment, targetComponent component2.Logging) (dataframe.DataFrame, error) {
	var implementationList []string
	var dataList []int
	var valueList []int
	for _, implementation := range t.implementations {
		data, values, err := experiment.Run(t.Client, targetComponent, t.Nameservers["sld"])
		if err != nil {
			return dataframe.DataFrame{}, nil
		}
		implementationList = append(implementationList, utils.Repeat(implementation.String(), len(data))...)
		dataList = append(dataList, data...)
		valueList = append(valueList, values...)
	}
	return dataframe.New(
		series.New(implementationList, series.String, "implementation"),
		series.New(dataList, series.Int, "data"),
		series.New(valueList, series.Int, "value"),
	), nil
}
