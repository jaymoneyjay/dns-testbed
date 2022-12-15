package testbed

import (
	"errors"
	"fmt"
	"os"
	"testbed/config"
	"text/template"
	"time"
)

type Zone struct {
	QName           string
	ZoneFileHost    string
	ZoneFileTarget  string
	DefaultZoneFile string
	*Container
	Implementation
}

func newZone(zoneConfig *config.Zone, templates string) *Zone {
	container := NewContainer(zoneConfig.ID, zoneConfig.Dir, zoneConfig.IP)
	impl := newBind(templates, container)
	return &Zone{
		QName:           zoneConfig.QName,
		ZoneFileHost:    zoneConfig.ZoneFileHost,
		ZoneFileTarget:  zoneConfig.ZoneFileTarget,
		DefaultZoneFile: zoneConfig.DefaultZoneFile,
		Container:       container,
		Implementation:  impl,
	}
}

func (z *Zone) Set(zoneFile string, reload bool) {
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
	if reload {
		z.reload()
	}
}

func (z *Zone) SetDefault(reload bool) {
	z.Set(z.DefaultZoneFile, reload)
}

func (z *Zone) SetDelay(duration time.Duration) {
	execResult, err := z.Exec([]string{"tc", "qdisc", "change", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", duration.Milliseconds())})
	if err != nil {
		panic(err)
	}
	if execResult.StdOut != "" {
		err = errors.New(fmt.Sprintf("could not set delay at %s: %s", z.QName, execResult.StdOut))
		panic(err)
	}
}
