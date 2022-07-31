package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"servicetitan-to-dataset/servicetitan"
	"strings"

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
	root.AddCommand(ListReportsCommand())

	return root
}

func askAuthQuestions(conf *servicetitan.ClientInfo) {
	askQuestion(conf, &conf.TenantID, "Tenant ID")
	askQuestion(conf, &conf.AppID, "App ID")
	askQuestion(conf, &conf.ClientID, "Client ID")
	askQuestion(conf, &conf.ClientSecret, "Client secret")
}

func askQuestion(conf *servicetitan.ClientInfo, attrRef *string, question string) {
	val, err := readValueFromInput(bufio.NewReader(os.Stdin), question)
	if err != nil {
		log.Fatal(err)
	}

	*attrRef = val
}

func readValueFromInput(reader *bufio.Reader, question string) (string, error) {
	fmt.Printf("Enter your serviceTitan %s: ", question)
	v, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Remove newline and carriage return for windows from value
	v = strings.TrimRight(v, "\n")
	v = strings.TrimRight(v, "\r")

	return v, nil
}

func readCredsFromEnv(conf *servicetitan.ClientInfo) {
	conf.TenantID = os.Getenv("SERVICE_TITAN_TENANT_ID")
	conf.AppID = os.Getenv("SERVICE_TITAN_APP_ID")
	conf.ClientID = os.Getenv("SERVICE_TITAN_CLIENT_ID")
	conf.ClientSecret = os.Getenv("SERVICE_TITAN_CLIENT_SECRET")
}
