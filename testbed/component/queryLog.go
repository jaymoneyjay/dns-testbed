package component

import (
	"os"
	"strings"
)

type QueryLog struct {
	logFile string
}

func NewLog(logFile string) *QueryLog {
	return &QueryLog{
		logFile: logFile,
	}
}

func (l *QueryLog) Clean() error {
	_, err := os.Create(l.logFile)
	if err != nil {
		return err
	}
	return nil
}

func (l *QueryLog) CountQueries() (int, error) {
	file, err := os.ReadFile(l.logFile)
	if err != nil {
		return 0, err
	}
	var cleanedByteSlice []uint8
	for _, byte := range file {
		if byte != 0 {
			cleanedByteSlice = append(cleanedByteSlice, byte)
		}
	}
	fileString := string(cleanedByteSlice)
	lines := strings.Split(fileString, "\n")
	return len(lines) - 1, nil
}
