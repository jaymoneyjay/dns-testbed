package cmd

import (
	"github.com/spf13/cobra"
	"testbed/config"
	"testbed/testbed"
)

var cmdZones = &cobra.Command{
	Use:   "zones [directory with zone files]",
	Short: "Set zone files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		testbedConfig, err := config.New().LoadTestbedConfig()
		if err != nil {
			return err
		}
		testbed.New(testbedConfig).SetZoneFiles(args[0])
		return nil
	},
}
