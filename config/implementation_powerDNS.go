package config

type powerDNS struct {
}

func newPowerDNS() powerDNS {
	return powerDNS{}
}

func (p powerDNS) configTarget() string {
	return "/etc/powerdns/recursor.conf"
}

func (p powerDNS) rootHintsTarget() string {
	return "/usr/share/dns/myroot.hints"
}

func (p powerDNS) logsTarget() string {
	//TODO
	return "/query.log"
}
