package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"time"
)

type Nameserver struct {
	id        string
	dockerCli *docker.Client
	implementation
}

func newNameserver(id string, client *docker.Client) *Nameserver {
	return &Nameserver{
		id:             id,
		dockerCli:      client,
		implementation: newBind("9.11.3", client),
	}
}

func (ns *Nameserver) Restart() {
	ns.implementation.restart(ns.id)
}

func (ns *Nameserver) FlushCache() {
	ns.implementation.flushCache(ns.id)
}

func (ns *Nameserver) SetZone(path string) {
	ns.dockerCli.WriteZoneFile(ns.id, path)
	ns.Restart()
}

func (ns *Nameserver) SetDelay(duration time.Duration) {
	execResult, err := ns.dockerCli.Exec(ns.id, []string{"tc", "qdisc", "change", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", duration.Milliseconds())})
	if err != nil {
		panic(err)
	}
	if execResult.StdOut != "" {
		err = errors.New(fmt.Sprintf("could not set delay at %s: %s", ns.id, execResult.StdOut))
		panic(err)
	}
}

func (ns *Nameserver) ID() string {
	return ns.id
}

func (ns *Nameserver) ReadQueryLog(minTimeout time.Duration) []byte {
	return ns.implementation.readQueryLog(ns.id, "nameserver", minTimeout)
}

func (ns *Nameserver) flushQueryLog() {
	ns.dockerCli.FlushLog(ns.id, "nameserver", "query.log")
}
