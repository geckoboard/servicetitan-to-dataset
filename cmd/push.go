package cmd

import (
	"context"
	"log"
	"os"
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/processor"
	"time"

	"github.com/spf13/cobra"
)

func PushDataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Fetch data from a serviceTitan and push to Geckoboard",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := loadAndValidateConfig(cmd.Flag("config").Value.String())
			if err != nil {
				log.Fatal(err)
			}

			if cfg.RefreshTimeSec == 0 {
				runAllEntries(context.Background(), cfg)
				log.Println("Completed pushing all entries")
				os.Exit(0)
			}

			// We use this instead of a ticker because we don't want
			// tickers to pile up as we need to wait 5mins between every
			// entry run
			for {
				runAllEntries(context.Background(), cfg)
				time.Sleep(time.Duration(cfg.RefreshTimeSec) * time.Second)
			}
		},
	}

	return cmd
}

func loadAndValidateConfig(path string) (*config.Config, error) {
	cfg, err := config.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return cfg, cfg.Validate()
}

func runAllEntries(ctx context.Context, cfg *config.Config) {
	for idx, ent := range cfg.Entries {
		proc := processor.New(cfg)

		log.Println("Processing entry...", idx)
		if err := proc.Process(ctx, ent); err != nil {
			log.Println("ERR: Unexpected error occurred", err)
		} else {
			log.Println("INF: Successfully processed and pushed")
		}

		// With every report request we make we have to wait another 5 minutes
		// to request the next report data as rate limits are 2 per 5 minutes.
		// Its a known issue that serviceTitan are working on resolving at some point
		log.Println("INF: Waiting 5 minutes for serviceTitan rate limit")
		time.Sleep(300 * time.Second)
	}
}
