package dns

import (
	"dns-testbed-go/docker"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

type unbound struct {
	kind      implementationKind
	version   string
	dockerCli *docker.Client
}

func newUnbound(version string, client *docker.Client) unbound {
	return unbound{
		kind:      Unbound,
		version:   version,
		dockerCli: client,
	}
}

func (u unbound) Kind() string {
	return u.kind.String()
}

func (u unbound) Version() string {
	return u.version
}

func (u unbound) restart(containerID string) {
	execResult, err := u.dockerCli.Exec(containerID, []string{"unbound-control", "reload"})
	if err != nil {
		panic(err)
	}
	reloadOK, err := regexp.MatchString("ok", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !reloadOK {
		err = errors.New(fmt.Sprintf("unbound cache could not be restarted successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (u unbound) flushCache(containerID string) {
	u.restart(containerID)
}

func (u unbound) readQueryLog(containerID, containerType string, minTimeout time.Duration) []byte {
	var lines []string
	numberOfCurrentLines := 0
	for true {
		time.Sleep(minTimeout)
		log := u.dockerCli.ReadLog(containerID, containerType, "query.log")
		lines = strings.Split(string(log), "\n")
		if len(lines) == numberOfCurrentLines {
			break
		}
		numberOfCurrentLines = len(lines)
	}
	lines = lines[0 : len(lines)-1]
	queries := u.filterQueries(lines)
	return []byte(strings.Join(queries, "\n"))
}

func (u unbound) filterQueries(lines []string) []string {
	var queries []string
	for _, line := range lines {
		matched, err := regexp.MatchString("(query:|reply:)", line)
		if err != nil {
			panic(err)
		}
		if matched {
			queries = append(queries, line)
		}
	}
	return queries
}
