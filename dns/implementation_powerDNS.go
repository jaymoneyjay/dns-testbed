package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"regexp"
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

func (p powerDNS) RestartExecution() execution {
	return execution{
		command: []string{"/etc/init.d/pdns-recursor", "restart"},
		responseVerification: func(response string) {
			matched, err := regexp.MatchString("Restarting PowerDNS Recursor \\.\\.\\.done", response)
			if err != nil {
				return
			}
			if err != nil {
				panic(err)
			}
			if !matched {
				err = errors.New(fmt.Sprintf("powerDNS cache could not be restarted successfully: %s", response))
				panic(err)
			}
		},
	}
}

func (p powerDNS) FlushCacheExecution() execution {
	return p.RestartExecution()
}
