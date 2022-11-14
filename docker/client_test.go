package docker

import (
	"fmt"
	"log"
	"regexp"
	"testing"
)

func TestShouldExecuteCommandOnContainer(t *testing.T) {
	execResult, err := NewClient().Exec("client", []string{"echo", "5"})
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)

		fmt.Print(execResult)
	}
}

func TestShouldDigBind11ForTargetCOM(t *testing.T) {
	testShouldDigZone("target.com.", "172.20.0.10", "172.20.0.7")
}

func TestShouldDigBind11ForInterNET(t *testing.T) {
	testShouldDigZone("inter.net.", "172.20.0.10", "172.20.0.8")
}

func TestShouldDigUnbound17ForTargetCOM(t *testing.T) {
	testShouldDigZone("target.com.", "172.20.0.11", "172.20.0.7")
}

func TestShouldDigUnbound17ForInterNET(t *testing.T) {
	testShouldDigZone("inter.net.", "172.20.0.11", "172.20.0.8")
}

func TestShouldDigUnbound16ForTargetCOM(t *testing.T) {
	testShouldDigZone("target.com.", "172.20.0.12", "172.20.0.7")
}

func TestShouldDigUnbound16ForInterNET(t *testing.T) {
	testShouldDigZone("inter.net.", "172.20.0.12", "172.20.0.8")
}

func TestShouldDigUnbound10ForTargetCOM(t *testing.T) {
	testShouldDigZone("target.com.", "172.20.0.13", "172.20.0.7")
}

func TestShouldDigUnbound10ForInterNET(t *testing.T) {
	testShouldDigZone("inter.net.", "172.20.0.13", "172.20.0.8")
}

func TestShouldDigPowerDNSForTargetCOM(t *testing.T) {
	testShouldDigZone("target.com.", "172.20.0.14", "172.20.0.7")
}

func TestShouldDigPowerDNSForInterNET(t *testing.T) {
	testShouldDigZone("inter.net.", "172.20.0.14", "172.20.0.8")
}

func TestShouldDigBind18ForTargetCOM(t *testing.T) {
	testShouldDigZone("target.com.", "172.20.0.15", "172.20.0.7")
}

func TestShouldDigBind18ForInterNET(t *testing.T) {
	testShouldDigZone("inter.net.", "172.20.0.15", "172.20.0.8")
}

func testShouldDigZone(zone, resolver, ip string) {
	execResult, err := NewClient().Exec("client", []string{"dig", zone, "A", fmt.Sprintf("@%s", resolver)})
	if err != nil {
		log.Fatal(err)
	}
	matched, err := regexp.MatchString(
		fmt.Sprintf(";; ANSWER SECTION:\n%s\t*[0-9]\t*IN\t*A\t%s", zone, ip),
		execResult.StdOut,
	)
	if err != nil {
		panic(err)
	}
	if !matched {
		log.Fatalf("could not dig %s for %s:\n%s", resolver, zone, execResult.StdOut)
	}
}
