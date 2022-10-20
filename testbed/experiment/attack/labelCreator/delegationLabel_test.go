package labelCreator

import (
	"log"
	"reflect"
	"testing"
)

func TestShouldIncrementDelegationLabel(t *testing.T) {
	delegationLabel := NewDelegationLabel()
	var actual []string
	for i := 0; i < 10; i++ {
		actual = append(actual, delegationLabel.Step())
	}

	expected := []string{
		"del1",
		"del2",
		"del3",
		"del4",
		"del5",
		"del6",
		"del7",
		"del8",
		"del9",
		"del10",
	}
	if !reflect.DeepEqual(actual, expected) {
		log.Fatalf("Actual: %s\nExpected: %s", actual, expected)
	}
}
