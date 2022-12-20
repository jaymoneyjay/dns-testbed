package config

import "path/filepath"

type unbound struct {
}

func newUnbound() unbound {
	return unbound{}
}

func (u unbound) configTarget() string {
	return filepath.Join("/usr", "local", "etc", "unbound", "unbound.conf")
}

func (u unbound) rootHintsTarget() string {
	return filepath.Join("/usr", "local", "etc", "unbound", "db.root")
}

func (u unbound) logsTarget() string {
	return filepath.Join("/usr", "local", "etc", "logs")
}

func (u unbound) startCommands() []string {
	return []string{
		"unbound-control-setup",
		"unbound-control start",
	}
}
