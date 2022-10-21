package testbed

import (
	"dns-testbed-go/testbed/component"
	"dns-testbed-go/testbed/docker"
	"dns-testbed-go/testbed/experiment"
	"dns-testbed-go/utils"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type Testbed struct {
	Nameservers     map[string][]*component.Nameserver
	Client          *component.Client
	Resolver        *component.Resolver
	dockerClient    *docker.Client
	implementations []component.Implementation
}

func NewTestbed() (*Testbed, error) {
	nameservers := map[string][]*component.Nameserver{
		"root": {component.NewNameserver("root", ".", "testbed/docker/buildContext/root")},
		"tld": {
			component.NewNameserver("com", "com.", "testbed/docker/buildContext/com"),
			component.NewNameserver("net", "net.", "testbed/docker/buildContext/net"),
		},
		"sld": {
			component.NewNameserver("target-com", "target.com.", "testbed/docker/buildContext/target-com"),
			component.NewNameserver("inter-net", "inter.net.", "testbed/docker/buildContext/inter-net"),
		},
	}
	for _, nsList := range nameservers {
		for _, ns := range nsList {
			err := ns.SetZoneFile("template.zone")
			if err != nil {
				return &Testbed{}, err
			}
		}
	}
	return &Testbed{
		Nameservers:     nameservers,
		Client:          component.NewClient("client"),
		Resolver:        component.NewResolver("resolver", "testbed/docker/buildContext/resolver"),
		dockerClient:    docker.NewClient(),
		implementations: []component.Implementation{component.Bind9},
	}, nil
}

func (t *Testbed) Start(implementation component.Implementation) error {
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

func (t *Testbed) Query(zone string) (docker.ExecResult, error) {
	return t.Client.Query(zone)
}

func (t *Testbed) Run(experiment *experiment.Experiment, targetComponent component.Logging) (dataframe.DataFrame, error) {
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

func (t *Testbed) CleanLogs() error {
	err := t.Resolver.CleanLog()
	if err != nil {
		return err
	}
	for _, tld := range t.Nameservers["tld"] {
		err := tld.CleanLog()
		if err != nil {
			return err
		}
	}
	for _, sld := range t.Nameservers["sld"] {
		err := sld.CleanLog()
		if err != nil {
			return err
		}
	}
	return nil
}
