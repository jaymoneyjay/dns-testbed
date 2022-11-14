package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
)

type Resolver struct {
	implementation
	*queryLog
	ip string
}

func newResolver(kind implementationKind, version, ip string, client *docker.Client) *Resolver {
	var implementation implementation
	switch kind {
	case Bind:
		implementation = newBind(version, client)
	case Unbound:
		implementation = newUnbound(version, client)
	case PowerDNS:
		implementation = newPowerDNS(version, client)
	default:
		panic(errors.New(fmt.Sprintf("implementation kind %s has no instantiation method yet.", kind)))
	}
	return &Resolver{
		ip:             ip,
		implementation: implementation,
		queryLog:       newQueryLog(fmt.Sprintf("%s-%s", implementation.Kind(), implementation.Version()), "resolver", client),
	}
}

func (r *Resolver) IP() string {
	return r.ip
}

func (r *Resolver) ID() string {
	return fmt.Sprintf("resolver-%s-%s", r.implementation.Kind(), r.implementation.Version())
}

func (r *Resolver) Restart() {
	exec := r.implementation.RestartExecution()
	execResult, err := r.dockerCli.Exec(r.ID(), exec.command)
	if err != nil {
		panic(err)
	}
	exec.responseVerification(execResult.StdOut)
}

func (r *Resolver) FlushCache() {
	exec := r.implementation.FlushCacheExecution()
	execResult, err := r.dockerCli.Exec(r.ID(), exec.command)
	if err != nil {
		panic(err)
	}
	exec.responseVerification(execResult.StdOut)
}
