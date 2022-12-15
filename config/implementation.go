package config

import (
	"errors"
	"fmt"
)

type implementation interface {
	configTarget() string
	rootHintsTarget() string
	logsTarget() string
	startCommands() []string
}

func newImplementation(impl string) implementation {
	switch impl {
	case "bind":
		return newBind()
	case "unbound":
		return newUnbound()
	case "powerdns":
		return newPowerDNS()
	default:
		err := errors.New(fmt.Sprintf("implementation %s not found while configuring", impl))
		panic(err)
	}
}
