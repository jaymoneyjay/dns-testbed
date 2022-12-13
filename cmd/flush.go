package cmd

import (
	"github.com/spf13/cobra"
	"testbed/config"
	"testbed/testbed"
)

var cmdFlush = &cobra.Command{
	Use:   "flush",
	Short: "Flush the cache of all resolvers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		testbedConfig, err := config.New().LoadTestbedConfig()
		if err != nil {
			return err
		}
		testbed.New(testbedConfig).Flush()
		return nil
	},
}
