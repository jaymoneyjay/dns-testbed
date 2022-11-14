package lab

type dataIterator struct {
	indexHue int
	dataHue  []string
	indexX   int
	dataX    []int
}

func NewDataIterator(dataHue []string, dataX []int) *dataIterator {
	return &dataIterator{
		dataHue:  dataHue,
		dataX:    dataX,
		indexX:   0,
		indexHue: 0,
	}
}

func (t *dataIterator) hasNextDataPoint() bool {
	return t.indexHue < len(t.dataHue) && t.indexX < len(t.dataX)
}

func (t *dataIterator) getNextDataPoint() (string, int) {
	nextHue := t.dataHue[t.indexHue]
	nextX := t.dataX[t.indexX]
	t.indexX += 1
	if t.indexX >= len(t.dataX) {
		t.indexX = 0
		t.indexHue += 1
	}
	return nextHue, nextX
}
