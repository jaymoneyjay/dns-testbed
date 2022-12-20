package config

import "path/filepath"

type bind struct {
}

func newBind() bind {
	return bind{}
}

func (b bind) configTarget() string {
	return filepath.Join("/etc", "bind", "named.conf.options")
}

func (b bind) rootHintsTarget() string {
	return filepath.Join("/usr", "share", "dns", "root.hints")
}

func (b bind) logsTarget() string {
	return filepath.Join("/etc", "logs")
}

func (b bind) startCommands() []string {
	return []string{"named -u bind -4 -d 2"}
}
