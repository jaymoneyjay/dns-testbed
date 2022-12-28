package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"testbed/config"
	"testbed/testbed"
)

var cmdInit = &cobra.Command{
	Use:     "init [testbed config]",
	Short:   "Initialize a dns testbed",
	Example: "testbed init validation/testbed-basic.yaml",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New()
		build := c.Load("build").(string)
		if _, err := os.Stat(build); !os.IsNotExist(err) {
			(&testbed.Testbed{Build: build}).Remove()
		}
		err := c.SetConfig(args[0], "testbed")
		if err != nil {
			return err
		}
		testbedConfig, err := c.LoadTestbedConfig()
		if err != nil {
			return err
		}
		testbed.Build(testbedConfig)
		fmt.Println("### Initialized testbed ###")
		fmt.Println(testbed.New(testbedConfig))
		return nil
	},
}
