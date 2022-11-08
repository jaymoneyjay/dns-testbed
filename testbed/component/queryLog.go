package component

import (
	"os"
	"strings"
	"time"
)

type QueryLog struct {
	LogFile string
}

func NewLog(logFile string) *QueryLog {
	return &QueryLog{
		LogFile: logFile,
	}
}

func (l *QueryLog) Clean() error {
	_, err := os.Create(l.LogFile)
	if err != nil {
		return err
	}
	return nil
}

func (l *QueryLog) CountQueries() (int, error) {
	lines, err := l.readCleanedLines()
	if err != nil {
		return 0, err
	}
	return len(lines) - 1, nil
}

func (l *QueryLog) GetQueryDuration(timeout time.Duration) (time.Duration, error) {
	var lines []string
	var err error
	/*numberOfCurrentLines := 0
	for true {
		time.Sleep(time.Second * 20)
		lines, err = l.readCleanedLines()
		if err != nil {
			return 0, err
		}
		if len(lines) == numberOfCurrentLines {
			break
		}
		numberOfCurrentLines = len(lines)
	}*/
	time.Sleep(timeout)
	lines, err = l.readCleanedLines()
	if err != nil {
		return 0, err
	}
	if len(lines) < 2 {
		return 0, nil
	}
	startTime, err := l.parseTimestamp(lines[0])
	if err != nil {
		return 0, err
	}
	endTime, err := l.parseTimestamp(lines[len(lines)-2])
	if err != nil {
		return 0, err
	}
	return endTime.Sub(startTime), nil
}

func (l *QueryLog) readCleanedLines() ([]string, error) {
	file, err := os.ReadFile(l.LogFile)
	if err != nil {
		return nil, err
	}
	var cleanedByteSlice []uint8
	for _, byte := range file {
		if byte != 0 {
			cleanedByteSlice = append(cleanedByteSlice, byte)
		}
	}
	fileString := string(cleanedByteSlice)
	return strings.Split(fileString, "\n"), nil
}

func (l *QueryLog) parseTimestamp(queryLine string) (time.Time, error) {
	timestamp := strings.Split(queryLine, " ")[0] + " " + strings.Split(queryLine, " ")[1]
	return time.Parse("02-Jan-2006 15:04:05.000", timestamp)
}
