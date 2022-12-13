/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"testbed/cmd"
)

func main() {
	cmd.Execute()
	/*testbedConfig, err := config.New().LoadTestbedConfig()
	if err != nil {
		panic(err)
	}
	testbed.New(testbedConfig).Start()
	testbed.New(testbedConfig).Stop()*/
}
