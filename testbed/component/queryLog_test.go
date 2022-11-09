package component

import (
	"log"
	"testing"
	"time"
)

func TestShouldGetQueryDuration(t *testing.T) {
	actual, _ := NewLog("test_query.log").GetQueryDuration()
	expected := time.Millisecond * 59464
	if actual != expected {
		log.Fatalf("Actual: %s\n Expected: %s", actual, expected)
	}
}
