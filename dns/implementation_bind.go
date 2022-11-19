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
	execResult, err := b.dockerCli.Exec(containerID, []string{"rndc", "flush"})
	if err != nil {
		panic(err)
	}
	flushedOK, err := regexp.MatchString("", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !flushedOK {
		err = errors.New(fmt.Sprintf("bind cache could not be flushed successfully: %s", execResult.StdOut))
		panic(err)
	}
	execResult, err = b.dockerCli.Exec(containerID, []string{"rndc", "reload"})
	reloadOK, err := regexp.MatchString("server reload successful", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !reloadOK {
		err = errors.New(fmt.Sprintf("bind cache could not be reloaded successfully: %s", execResult.StdOut))
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
