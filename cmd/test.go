package cmd

import (
	"github.com/mono83/mtm/prometheus"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"dry", "dry-run"},
	Short:   "Perform dry run to test configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := genericReader()()
		if err != nil {
			return err
		}

		prometheus.WriteGaugeTo(cmd.OutOrStdout(), data)
		return nil
	},
}

func init() {
	genericInject(testCmd, false, false)
}
