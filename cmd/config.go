package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"servicetitan-to-dataset/config"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

func ConfigCommand() *cobra.Command {
	var (
		generate bool
		validate bool
	)

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Generate and validate a config",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err        error
				configPath = cmd.Flag("config").Value.String()
			)
			switch {
			case generate:
				err = buildExampleConfig(configPath)
			case validate:
				cfg, err := config.LoadFile(configPath)
				if err != nil {
					log.Fatal(err)
				}

				if err := cfg.Validate(); err != nil {
					log.Fatal(err)
				}

				log.Println("Config all valid...")
			default:
				err = errors.New("missing --generate or --validate switch")
			}

			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVar(&generate, "generate", false, "Generate a template config")
	cmd.Flags().BoolVar(&validate, "validate", false, "Validate a config")

	return cmd
}

func buildExampleConfig(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("%s already exists... please rename or delete or use --config to specify a different output", filename)
	}

	cfg := config.Config{
		ServiceTitan: config.ServiceTitan{
			AppID:        "{{ENV_APPID}}",
			TenantID:     "your-tenant-id",
			ClientID:     "your-client-id",
			ClientSecret: "your-client-secret",
		},
		Geckoboard: config.Geckoboard{
			APIKey: "apikey1234",
		},
		RefreshTimeSec: 60,
		Entries: config.Entries{
			{
				Report: config.Report{
					ID:         "123",
					CategoryID: "category-a",
					Parameters: []config.Parameter{
						{
							Name:  "From",
							Value: "2021-10-13",
						},
						{
							Name:  "To",
							Value: "2021-10-14",
						},
					},
				},
				Dataset: config.Dataset{
					Name: "my-dataset-name",
					Type: "replace",
				},
			},
			{
				Report: config.Report{
					ID:         "345",
					CategoryID: "category-b",
				},
				Dataset: config.Dataset{
					Name: "revenue-income",
					Type: "append",
				},
			},
		},
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	return yaml.NewEncoder(f).Encode(cfg)
}
