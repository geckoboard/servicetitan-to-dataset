package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = ""

func Setup() *cobra.Command {
	var configPath string

	root := &cobra.Command{
		Use:   "servicetitan-to-dataset",
		Short: `Push your service titan reports to your Geckoboard dataset`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Hidden: true,
	}

	root.Run = func(cmd *cobra.Command, args []string) {
		curr, _, _ := root.Find(os.Args[1:])

		// Default to help if no commands present
		if curr.Use == root.Use {
			root.SetArgs([]string{"-h"})
			root.Execute()
		}
	}

	root.PersistentFlags().StringVar(&configPath, "config", "config.yml", "Path to the config file")

	root.AddCommand(VersionCommand())
	root.AddCommand(ConfigCommand())
	root.AddCommand(ReportsCommand())
	root.AddCommand(PushDataCommand())

	return root
}
