package docker

import (
	"fmt"
	"log"
	"testing"
)

func TestShouldExecuteCommandOnContainer(t *testing.T) {
	client, _ := NewClient()
	execResult, err := client.Exec("client", []string{"echo", "5"})
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(execResult)
}
