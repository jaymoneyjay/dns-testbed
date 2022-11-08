package component

import (
	"dns-testbed-go/testbed/docker"
	"fmt"
	"os"
	"path/filepath"
)

type Nameserver struct {
	*QueryLog
	*Container
	zone      string
	buildPath string
	zonePath  string
}

func attachNameserver(containerID, zone, buildPath string) (*Nameserver, error) {
	container, err := NewContainer(containerID)
	if err != nil {
		return nil, err
	}
	return &Nameserver{
		Container: container,
		zone:      zone,
		buildPath: buildPath,
		zonePath:  filepath.Join(buildPath, "zones"),
		QueryLog:  NewLog(filepath.Join(buildPath, "logs/query.log")),
	}, nil
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

func (ns *Nameserver) SetZoneFile(zoneFileID string) error {
	localTemplate := fmt.Sprintf(`zone "%s" {
		type master;
		file "/etc/zones/%s";
	};
	`, ns.zone, zoneFileID)
	err := os.WriteFile(filepath.Join(ns.buildPath, "named.conf.local"), []byte(localTemplate), 0644)
	if err != nil {
		return err
	}
	return ns.restartBind9()
}

func (ns *Nameserver) Restart() error {
	return ns.restartBind9()
}

func (ns *Nameserver) SetDelay(delay int) (docker.ExecResult, error) {
	return ns.Exec([]string{"tc", "qdisc", "change", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", delay)})
}

func AttachRoot() (*Nameserver, error) {
	return attachNameserver("root", ".", "testbed/docker/buildContext/nameserver/root")
}

func AttachCOM() (*Nameserver, error) {
	return attachNameserver("com", "com.", "testbed/docker/buildContext/nameserver/com")
}

func AttachNET() (*Nameserver, error) {
	return attachNameserver("net", "net.", "testbed/docker/buildContext/nameserver/net")
}

func AttachTarget() (*Nameserver, error) {
	return attachNameserver("target-com", "target.com.", "testbed/docker/buildContext/nameserver/target-com")
}

func AttachInter() (*Nameserver, error) {
	return attachNameserver("inter-net", "inter.net.", "testbed/docker/buildContext/nameserver/inter-net")
}
