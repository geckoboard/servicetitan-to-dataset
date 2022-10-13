package cmd

import (
	"servicetitan-to-dataset/cmd/report"

	"github.com/spf13/cobra"
)

func ReportsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reports",
		Short: "List reports and print specific report parameters required for the config",
	}

	cmd.AddCommand(report.ListCommand())

	return cmd
}
