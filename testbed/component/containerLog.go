package component

import (
	"os"
	"strings"
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
	file, err := os.Create(l.generalLog)
	file.Close()
	if err != nil {
		return err
	}
	file, err = os.Create(l.queryLog)
	file.Close()
	return err
}

func (l *containerLog) CountQueries() (int, error) {
	file, err := os.ReadFile(l.queryLog)
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
