package cmd

import (
	"github.com/spf13/cobra"
	"testbed/config"
	"testbed/testbed"
)

var cmdInit = &cobra.Command{
	Use:   "init [testbed config, zone files]",
	Short: "Initialize a dns testbed",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New()
		if previousTestbedConfig, err := c.LoadTestbedConfig(); err == nil {
			testbed.New(previousTestbedConfig).Remove()
		}
		err := c.SetConfig(args[0])
		if err != nil {
			return err
		}
		testbedConfig, err := c.LoadTestbedConfig()
		if err != nil {
			return err
		}
		testbed.Build(testbedConfig, args[1])
		return nil
	},
}
