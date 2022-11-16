package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type bind struct {
	dockerCli *docker.Client
	kind      implementationKind
	version   string
}

func newBind(version string, client *docker.Client) bind {
	return bind{
		kind:      Bind,
		version:   version,
		dockerCli: client,
	}
}

func (b bind) Kind() string {
	return b.kind.String()
}

func (b bind) Version() string {
	return b.version
}

func (b bind) restart(containerID string) {
	execResult, err := b.dockerCli.Exec(containerID, []string{"service", "named", "restart"})
	if err != nil {
		panic(err)
	}
	patternStopOK := "(\\* Stopping)[^*]*done"
	patternStartOK := "(\\* Starting)[^*]*done"
	stoppedOK, err := regexp.MatchString(patternStopOK, execResult.StdOut)
	if err != nil {
		panic(err)
	}
	startedOK, err := regexp.MatchString(patternStartOK, execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !(stoppedOK && startedOK) {
		err = errors.New(fmt.Sprintf("bind cache could not be restarted successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (b bind) flushCache(containerID string) {
	b.restart(containerID)
}

func (b bind) readQueryLog(containerID, containerType string, minTimeout time.Duration) []byte {
	var lines []string
	numberOfCurrentLines := 0
	for true {
		time.Sleep(minTimeout)
		log := b.dockerCli.ReadLog(containerID, containerType, "query.log")
		lines = strings.Split(string(log), "\n")
		if len(lines) == numberOfCurrentLines {
			break
		}
		numberOfCurrentLines = len(lines)
	}
	queries := b.filterQueries(lines)
	return []byte(strings.Join(queries, "\n"))
}

func (b bind) filterQueries(lines []string) []string {
	var queries []string
	for _, line := range lines {
		matched, err := regexp.MatchString("(query:)", line)
		if err != nil {
			panic(err)
		}
		if matched {
			queries = append(queries, line)
		}
	}
	return queries
}
