package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"testbed/config"
	"testbed/experiment"
	"testbed/testbed"
)

var runAll bool

var cmdRun = &cobra.Command{
	Use:     "run [experiment config]",
	Short:   "Run an experiment according to the specified configuration",
	Example: "run validation/experiments/subquery+CNAME.yaml ",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := config.New()
		testbedConfig, err := c.LoadTestbedConfig()
		if err != nil {
			return err
		}
		t := testbed.New(testbedConfig)
		configStat, err := os.Stat(args[0])
		if err != nil {
			return nil
		}
		if runAll {
			if !configStat.IsDir() {
				fmt.Println("must provide directory when running with -a")
			}
			configEntries, err := os.ReadDir(args[0])
			if err != nil {
				return err
			}
			for _, entry := range configEntries {
				configPath := filepath.Join(args[0], entry.Name())
				if err := runExperiment(configPath, c, t); err != nil {
					fmt.Printf("cannot run experiment with configuration %s", configPath)
				}
			}
		}
		experimentConfig, err := c.LoadExperimentConfig(args[0])
		if err != nil {
			return err
		}
		if err := experiment.New(experimentConfig).Run(t); err != nil {
			return err
		}
		fmt.Printf("results written to %s\n", experimentConfig.Dest)
		return nil
	},
}

func runExperiment(experimentConfigPath string, c *config.Config, t *testbed.Testbed) error {
	experimentConfig, err := c.LoadExperimentConfig(experimentConfigPath)
	if err != nil {
		return err
	}
	if err := experiment.New(experimentConfig).Run(t); err != nil {
		return err
	}
	fmt.Printf("results written to %s\n", experimentConfig.Dest)
	return nil
}

func init() {
	cmdRun.Flags().BoolVarP(&runAll, "runAll", "a", false, "Run all experiments at configuration location")
}
