package dns

import (
	"dns-testbed-go/docker"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

type unbound struct {
	*queryLog
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

func (u unbound) RestartExecution() execution {
	return execution{
		command: []string{"unbound-control", "reload"},
		responseVerification: func(response string) {
			reloadOK, err := regexp.MatchString("ok", response)
			if err != nil {
				panic(err)
			}
			if !reloadOK {
				err = errors.New(fmt.Sprintf("unbound cache could not be restarted successfully: %s", response))
				panic(err)
			}
		},
	}

}

func (u unbound) FlushCacheExecution() execution {
	return u.RestartExecution()
}
