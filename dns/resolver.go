package dns

import (
	"dns-testbed-go/docker"
	"errors"
	"fmt"
	"time"
)

type Resolver struct {
	implementation
	ip        string
	dockerCli *docker.Client
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
		dockerCli:      client,
	}
}

func (r *Resolver) IP() string {
	return r.ip
}

func (r *Resolver) ID() string {
	return fmt.Sprintf("resolver-%s-%s", r.implementation.Kind(), r.implementation.Version())
}

func (r *Resolver) Restart() {
	r.implementation.restart(r.ID())
}

func (r *Resolver) FlushCache() {
	r.implementation.flushCache(r.ID())
}

func (r *Resolver) ReadQueryLog(minTimeout time.Duration) []byte {
	return r.implementation.readQueryLog(r.ID(), "resolver", minTimeout)
}

func (r *Resolver) flushQueryLog() {
	r.dockerCli.FlushLog(r.ID(), "resolver", "query.log")
}
