package config

import "path/filepath"

type powerDNS struct {
}

func newPowerDNS() powerDNS {
	return powerDNS{}
}

func (p powerDNS) configTarget() string {
	return filepath.Join("/etc", "powerdns", "recursor.conf")
}

func (p powerDNS) rootHintsTarget() string {
	return filepath.Join("/usr", "share", "dns", "myroot.hints")
}

func (p powerDNS) logsTarget() string {
	//TODO
	return filepath.Join("/")
}

func (p powerDNS) startCommands() []string {
	return []string{"/etc/init.d/pdns-recursor start"}
}
