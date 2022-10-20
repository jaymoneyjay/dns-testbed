package component

import (
	"bufio"
	"os"
)

type log struct {
	path string
}

func newLog(path string) *log {
	return &log{
		path: path,
	}
}

func (l *log) Clean() error {
	return os.WriteFile(l.path, []byte(""), 0666)
}

func (l *log) CountQueries() (int, error) {
	logs, err := os.Open(l.path)
	if err != nil {
		return 0, err
	}
	logsScanner := bufio.NewScanner(logs)
	queryCount := 0
	for logsScanner.Scan() {
		queryCount++
	}
	return queryCount, nil
}
