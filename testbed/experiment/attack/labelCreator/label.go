package labelCreator

import (
	"fmt"
)

type label struct {
	prefix  string
	counter int
}

func newLabel(prefix string) *label {
	return &label{prefix: prefix, counter: 0}
}

func (l *label) Step() string {
	l.counter++
	return l.String()
}

func (l *label) String() string {
	return fmt.Sprintf("%s%d", l.prefix, l.counter)
}
