package templates

type Zone struct {
	QName    string
	ZoneFile string
}

type NameServer struct {
	Zones []*Zone
}

type Resolver struct {
	QMin string
}
