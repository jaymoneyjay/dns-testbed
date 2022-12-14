package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"testbed/config"
	"testbed/testbed"
)

var cmdInit = &cobra.Command{
	Use:   "init [testbed config, zone files]",
	Short: "Initialize a dns testbed",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New()
		build := c.Load("build").(string)
		if _, err := os.Stat(build); !os.IsNotExist(err) {
			(&testbed.Testbed{Build: build}).Remove()
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
