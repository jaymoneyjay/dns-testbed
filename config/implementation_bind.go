package config

type bind struct {
}

func newBind() bind {
	return bind{}
}

func (b bind) configTarget() string {
	return "/etc/bind/named.conf.options"
}

func (b bind) rootHintsTarget() string {
	return "/usr/share/dns/root.hints"
}

func (b bind) logsTarget() string {
	return "/etc/logs/query.log"
}

func (b bind) startCommands() []string {
	return []string{"named -u bind -4 -d 2"}
}
