package dns

import (
	"dns-testbed-go/docker"
	"strings"
	"time"
)

type queryLog struct {
	id        string
	kind      string
	dockerCli *docker.Client
}

func newQueryLog(id, kind string, client *docker.Client) *queryLog {
	return &queryLog{
		id:        id,
		dockerCli: client,
		kind:      kind,
	}
}

func (q *queryLog) ReadQueryLog(minTimeout time.Duration) []byte {
	var queryLog []byte
	numberOfCurrentLines := 0
	for true {
		time.Sleep(minTimeout)
		queryLog = q.dockerCli.ReadLog(q.id, q.kind, "query.log")
		lines := strings.Split(string(queryLog), "\n")
		if len(lines) == numberOfCurrentLines {
			break
		}
		numberOfCurrentLines = len(lines)
	}
	return queryLog
}

func (q *queryLog) flushQueryLog() {
	q.dockerCli.FlushLog(q.id, q.kind, "query.log")
}
