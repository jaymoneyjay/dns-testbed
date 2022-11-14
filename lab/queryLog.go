package lab

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type queryLog struct {
	targetLog   []byte
	interLog    []byte
	resolverLog []byte
	hue         string
	x           int
}

func newQueryLog(targetLog, interLog, resolverLog []byte, hue string, x int) *queryLog {
	return &queryLog{
		targetLog:   targetLog,
		interLog:    interLog,
		resolverLog: resolverLog,
		hue:         hue,
		x:           x,
	}
}

func (q *queryLog) save(logDir string) {
	perm := fs.ModePerm
	currentLog := filepath.Join(logDir, fmt.Sprintf("%s-%d", q.hue, q.x))
	err := os.Mkdir(currentLog, perm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(
		filepath.Join(currentLog, "target.log"),
		q.targetLog,
		perm,
	)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(
		filepath.Join(currentLog, "inter.log"),
		q.interLog,
		perm,
	)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(
		filepath.Join(currentLog, "resolver.log"),
		q.resolverLog,
		perm,
	)
	if err != nil {
		panic(err)
	}
}
