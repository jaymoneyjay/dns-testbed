package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"regexp"
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

func (b bind) RestartExecution() execution {
	return execution{
		command: []string{"service", "bind9", "restart"},
		responseVerification: func(response string) {
			patternStopOK := "(\\* Stopping)[^*]*done"
			patternStartOK := "(\\* Starting)[^*]*done"
			stoppedOK, err := regexp.MatchString(patternStopOK, response)
			if err != nil {
				panic(err)
			}
			startedOK, err := regexp.MatchString(patternStartOK, response)
			if err != nil {
				panic(err)
			}
			if !(stoppedOK && startedOK) {
				err = errors.New(fmt.Sprintf("bind cache could not be restarted successfully: %s", response))
				panic(err)
			}
		},
	}
}

func (b bind) FlushCacheExecution() execution {
	return b.RestartExecution()
}
