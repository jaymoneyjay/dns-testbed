package labelCreator

import (
	"log"
	"reflect"
	"testing"
)

func TestShouldIncrementChainLabel(t *testing.T) {
	chainLabel := NewChainLabel()
	var actual []string
	for i := 0; i < 10; i++ {
		actual = append(actual, chainLabel.Step())
	}

	expected := []string{
		"chain1",
		"chain2",
		"chain3",
		"chain4",
		"chain5",
		"chain6",
		"chain7",
		"chain8",
		"chain9",
		"chain10",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}
