package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"time"
)

type Nameserver struct {
	*queryLog
	ID        string
	dockerCli *docker.Client
	implementation
}

func newNameserver(id string, client *docker.Client) *Nameserver {
	return &Nameserver{
		ID:             id,
		dockerCli:      client,
		queryLog:       newQueryLog(id, "nameserver", client),
		implementation: newBind("9.11.3", client),
	}
}

func (ns *Nameserver) Restart() {
	exec := ns.implementation.RestartExecution()
	execResult, err := ns.dockerCli.Exec(ns.id, exec.command)
	if err != nil {
		panic(err)
	}
	exec.responseVerification(execResult.StdOut)
}

func (ns *Nameserver) SetZone(path string) {
	ns.dockerCli.WriteZoneFile(ns.ID, path)
	ns.Restart()
}

func (ns *Nameserver) SetDelay(duration time.Duration) {
	execResult, err := ns.dockerCli.Exec(ns.ID, []string{"tc", "qdisc", "change", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", duration.Milliseconds())})
	if err != nil {
		panic(err)
	}
	if execResult.StdOut != "" {
		err = errors.New(fmt.Sprintf("could not set delay at %s: %s", ns.ID, execResult.StdOut))
		panic(err)
	}
}
