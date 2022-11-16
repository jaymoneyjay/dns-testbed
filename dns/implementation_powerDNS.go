package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type powerDNS struct {
	kind      implementationKind
	version   string
	dockerCli *docker.Client
}

func newPowerDNS(version string, client *docker.Client) powerDNS {
	return powerDNS{
		kind:      PowerDNS,
		version:   version,
		dockerCli: client,
	}
}

func (p powerDNS) Kind() string {
	return p.kind.String()
}

func (p powerDNS) Version() string {
	return p.version
}

func (p powerDNS) restart(containerID string) {
	execResult, err := p.dockerCli.Exec(containerID, []string{"/etc/init.d/pdns-recursor", "restart"})
	if err != nil {
		panic(err)
	}
	matched, err := regexp.MatchString("Restarting PowerDNS Recursor \\.\\.\\.done", execResult.StdOut)
	if err != nil {
		return
	}
	if err != nil {
		panic(err)
	}
	if !matched {
		err = errors.New(fmt.Sprintf("powerDNS cache could not be restarted successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (p powerDNS) flushCache(containerID string) {
	p.restart(containerID)
}

func (p powerDNS) readQueryLog(containerID, containerType string, minTimeout time.Duration) []byte {
	//TODO implement
	var lines []string
	numberOfCurrentLines := 0
	for true {
		time.Sleep(minTimeout)
		log := p.dockerCli.ReadLog(containerID, containerType, "query.log")
		lines = strings.Split(string(log), "\n")
		if len(lines) == numberOfCurrentLines {
			break
		}
		numberOfCurrentLines = len(lines)
	}
	queries := p.filterQueries(lines)
	return []byte(strings.Join(queries, "\n"))
}

func (p powerDNS) filterQueries(lines []string) []string {
	//TODO implement
	return lines
}
