package lab

import (
	"log"
	"os"
	"testing"
)

func TestComputeQueryDuration(t *testing.T) {
	testExperiment := newTimingExperiment("test", "testzone")
	file, err := os.ReadFile("test_querylog")
	if err != nil {
		panic(err)
	}
	actual, err := testExperiment.computeQueryDuration(file)
	if err != nil {
		panic(err)
	}
	if actual.Seconds() != 4.0 {
		log.Fatalf("Actual: %d, Expected: %d", actual, 4)
	}
}
