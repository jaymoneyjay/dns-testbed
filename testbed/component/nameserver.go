package component

import (
	"fmt"
	"os"
	"path/filepath"
)

type Nameserver struct {
	*Container
	zone      string
	buildPath string
	zonePath  string
	log       *log
}

func NewNameserver(containerID, zone, buildPath string) *Nameserver {
	return &Nameserver{
		Container: newContainer(containerID),
		zone:      zone,
		buildPath: buildPath,
		zonePath:  filepath.Join(buildPath, "zones"),
		log:       newLog(filepath.Join(buildPath, "logs/query.log"), filepath.Join(buildPath, "logs/general.log")),
	}
}

func (ns *Nameserver) WriteZone(zoneFragment, zoneFileID string) error {
	err := ns.SetZoneFile(zoneFileID)
	if err != nil {
		return err
	}
	template, err := os.ReadFile(filepath.Join(ns.zonePath, "template.zone"))
	if err != nil {
		return err
	}
	zoneData := fmt.Sprintf("%s\n%s", template, zoneFragment)
	return os.WriteFile(filepath.Join(ns.zonePath, zoneFileID), []byte(zoneData), 0666)
}

func (ns *Nameserver) GetZone() string {
	return ns.zone
}

func (ns Nameserver) CleanLog() error {
	return ns.log.Clean()
}

func (ns Nameserver) CountQueries() (int, error) {
	return ns.log.CountQueries()
}

func (ns *Nameserver) SetZoneFile(zoneFileID string) error {
	localTemplate := fmt.Sprintf(`zone "%s" {
		type master;
		file "/etc/zones/%s";
	};
	`, ns.zone, zoneFileID)
	err := os.WriteFile(filepath.Join(ns.buildPath, "bind9", "named.conf.local"), []byte(localTemplate), 0644)
	if err != nil {
		return err
	}
	return ns.startBind9()
}
