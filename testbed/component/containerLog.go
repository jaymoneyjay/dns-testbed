package component

import (
	"bufio"
	"os"
)

type containerLog struct {
	queryLog   string
	generalLog string
}

func newLog(queryLogPath, generalLogPath string) *containerLog {
	return &containerLog{
		queryLog:   queryLogPath,
		generalLog: generalLogPath,
	}
}

func (l *containerLog) Clean() error {
	_, err := os.Create(l.generalLog)
	if err != nil {
		return err
	}
	_, err = os.Create(l.queryLog)
	return err
}

func (l *containerLog) CountQueries() (int, error) {
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
