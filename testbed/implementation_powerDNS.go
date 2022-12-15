package testbed

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

type powerDNS struct {
	templatesDir string
	container    *Container
}

func newPowerDNS(templatesDir string, container *Container) powerDNS {
	return powerDNS{
		templatesDir: templatesDir,
		container:    container,
	}
}

func (p powerDNS) reload() {
	execResult, err := p.container.Exec([]string{"/etc/init.d/pdns-recursor", "restart"})
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
		err = errors.New(fmt.Sprintf("powerDNS could not be reloaded successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (p powerDNS) start() {
	execResult, err := p.container.Exec([]string{"/etc/init.d/pdns-recursor", "start"})
	if err != nil {
		panic(err)
	}
	matched, err := regexp.MatchString("done", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !matched {
		err = errors.New(fmt.Sprintf("powerDNS could not be started successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (p powerDNS) flushCache() {
	p.reload()
}

func (p powerDNS) filterQueries(queryLog []byte) []byte {
	//TODO implement
	return queryLog
}

func (p powerDNS) SetConfig(qmin, reload bool) {
	tmpl, err := template.ParseFiles(filepath.Join(p.templatesDir, "resolver-powerdns.conf"))
	if err != nil {
		panic(err)
	}
	dest, err := os.Create(filepath.Join(p.container.dir, "powerdns.conf"))
	if err != nil {
		panic(err)
	}
	var param string
	if qmin {
		param = "yes"
	} else {
		param = "no"
	}
	err = tmpl.Execute(dest, param)
	if err != nil {
		panic(err)
	}
	if reload {
		p.reload()
	}
}
