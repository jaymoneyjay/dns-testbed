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
	implementations []component.Implementation
}

func NewTestbed() (*Testbed, error) {
	root, err := component.NewNameserver("root", ".", "testbed/docker/buildContext/root")
	if err != nil {
		return nil, err
	}
	com, err := component.NewNameserver("com", "com.", "testbed/docker/buildContext/com")
	if err != nil {
		return nil, err
	}
	net, err := component.NewNameserver("net", "net.", "testbed/docker/buildContext/net")
	if err != nil {
		return nil, err
	}
	target, err := component.NewNameserver("target-com", "target.com.", "testbed/docker/buildContext/target-com")
	if err != nil {
		return nil, err
	}
	inter, err := component.NewNameserver("inter-net", "inter.net.", "testbed/docker/buildContext/inter-net")
	if err != nil {
		return nil, err
	}
	nameservers := map[string][]*component.Nameserver{
		"root": {root},
		"tld": {
			com,
			net,
		},
		"sld": {
			target,
			inter,
		},
	}
	for _, nsList := range nameservers {
		for _, ns := range nsList {
			err := ns.Start()
			if err != nil {
				return &Testbed{}, err
			}
		}
	}
	client, err := component.NewClient("client")
	if err != nil {
		return nil, err
	}
	resolver, err := component.NewResolver("resolver", "testbed/docker/buildContext/resolver")
	if err != nil {
		return nil, err
	}
	return &Testbed{
		Nameservers:     nameservers,
		Client:          client,
		Resolver:        resolver,
		implementations: []component.Implementation{component.Bind9},
	}, nil
}

func (t *Testbed) Start(implementation component.Implementation) error {
	return t.Resolver.Start(implementation)
}

func (t *Testbed) Stop(implementation component.Implementation) error {
	return t.Resolver.Stop(implementation)
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
	for _, nsList := range t.Nameservers {
		for _, nameserver := range nsList {
			err := nameserver.CleanLog()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
