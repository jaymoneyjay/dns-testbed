package testbed

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testbed/testbed/templates"
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
	_, err := p.container.Exec([]string{"/etc/init.d/pdns-recursor", "stop"})
	if err != nil {
		panic(err)
	}
	execResult, err := p.container.Exec([]string{"pdns_recursor", "--config-dir=/etc/powerdns/config/", "--daemon"})
	if err != nil {
		panic(err)
	}
	matched, err := regexp.MatchString("Exception:", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if matched {
		err = errors.New(fmt.Sprintf("powerDNS could not be started successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (p powerDNS) start() {
	execResult, err := p.container.Exec([]string{"pdns_recursor", "--config-dir=/etc/powerdns/config/", "--daemon"})
	if err != nil {
		panic(err)
	}
	matched, err := regexp.MatchString("Exception:", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if matched {
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
	tmpl, err := template.ParseFiles(filepath.Join(p.templatesDir, "recursor.conf"))
	if err != nil {
		panic(err)
	}
	dest, err := os.Create(filepath.Join(p.container.Config, "recursor.conf"))
	if err != nil {
		panic(err)
	}
	options := &templates.Resolver{
		QMin: "no",
	}
	if qmin {
		options.QMin = "yes"
	}
	err = tmpl.Execute(dest, options)
	if err != nil {
		panic(err)
	}
	if reload {
		p.reload()
	}
}
