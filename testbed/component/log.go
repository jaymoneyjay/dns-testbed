package component

import (
	"bufio"
	"os"
)

type log struct {
	queryLog   string
	generalLog string
}

func newLog(queryLogPath, generalLogPath string) *log {
	return &log{
		queryLog:   queryLogPath,
		generalLog: generalLogPath,
	}
}

func (l *log) Clean() error {
	err := os.WriteFile(l.generalLog, []byte(""), 0666)
	if err != nil {
		return err
	}
	return os.WriteFile(l.queryLog, []byte(""), 0666)
}

func (l *log) CountQueries() (int, error) {
	logs, err := os.Open(l.queryLog)
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
