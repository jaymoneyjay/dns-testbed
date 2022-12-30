package testbed

import (
	"os"
	"testbed/config"
	"text/template"
)

type Zone struct {
	NS              *Nameserver
	QName           string
	ZoneFileHost    string
	ZoneFileTarget  string
	DefaultZoneFile string
}

func newZone(zoneConfig *config.Zone) *Zone {
	return &Zone{
		QName:           zoneConfig.QName,
		ZoneFileHost:    zoneConfig.ZoneFileHost,
		ZoneFileTarget:  zoneConfig.ZoneFileTarget,
		DefaultZoneFile: zoneConfig.DefaultZoneFile,
	}
}

func (z *Zone) set(zoneFile string) {
	tmpl, err := template.ParseFiles(zoneFile)
	if err != nil {
		panic(err)
	}
	dest, err := os.Create(z.ZoneFileHost)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(dest, nil)
	if err != nil {
		return
	}
}

func (z *Zone) setDefault() {
	z.set(z.DefaultZoneFile)
}
