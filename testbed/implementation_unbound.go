package testbed

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type unbound struct {
	templatesDir string
	container    *Container
}

func newUnbound(templatesDir string, container *Container) unbound {
	return unbound{
		templatesDir: templatesDir,
		container:    container,
	}
}

func (u unbound) reload() {
	execResult, err := u.container.Exec([]string{"unbound-control", "reload"})
	if err != nil {
		panic(err)
	}
	reloadOK, err := regexp.MatchString("ok", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !reloadOK {
		err = errors.New(fmt.Sprintf("unbound cache could not be restarted successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (u unbound) start() {
	execResult, err := u.container.Exec([]string{"unbound-control-setup"})
	if err != nil {
		panic(err)
	}
	matched, err := regexp.MatchString("Setup success", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !matched {
		err = errors.New(fmt.Sprintf("unbound setup not executed successfully: %s", execResult.StdOut))
		panic(err)
	}
	execResult, err = u.container.Exec([]string{"unbound-control", "start"})
	if err != nil {
		panic(err)
	}
	matched, err = regexp.MatchString("", execResult.StdOut)
	if err != nil {
		panic(err)
	}
	if !matched {
		err = errors.New(fmt.Sprintf("unbound start not executed successfully: %s", execResult.StdOut))
		panic(err)
	}
}

func (u unbound) flushCache() {
	u.reload()
}

func (u unbound) filterQueries(queryLog []byte) []byte {
	var queries []string
	lines := strings.Split(string(queryLog), "\n")
	for _, line := range lines {
		matched, err := regexp.MatchString("(query:|reply:)", line)
		if err != nil {
			panic(err)
		}
		if matched {
			queries = append(queries, line)
		}
	}
	return []byte(strings.Join(queries, "\n"))
}

func (u unbound) setConfig(qmin, reload bool) {
	tmpl, err := template.ParseFiles(filepath.Join(u.templatesDir, "resolver-unbound.conf"))
	if err != nil {
		panic(err)
	}
	dest, err := os.Create(filepath.Join(u.container.dir, "unbound.conf"))
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
		u.reload()
	}
}
