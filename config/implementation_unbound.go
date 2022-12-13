package config

type unbound struct {
}

func newUnbound() unbound {
	return unbound{}
}

func (u unbound) configTarget() string {
	return "/usr/local/etc/unbound/unbound.conf"
}

func (u unbound) rootHintsTarget() string {
	return "/usr/local/etc/unbound/db.root"
}

func (u unbound) logsTarget() string {
	return "/usr/local/etc/logs/query.log"
}
