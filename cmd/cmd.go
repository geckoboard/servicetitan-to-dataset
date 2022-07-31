package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = ""

func Setup() *cobra.Command {
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

	root.AddCommand(VersionCommand())

	return root
}
