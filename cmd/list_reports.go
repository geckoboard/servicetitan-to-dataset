package cmd

import (
	"context"
	"log"
	"os"
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/servicetitan"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type categoryReportEntry struct {
	reports  []servicetitan.Report
	category servicetitan.Category
}

func ListReportsCommand() *cobra.Command {
	var (
		credsFromEnv bool
		filter       string
	)

	cmd := &cobra.Command{
		Use:   "list-reports",
		Short: "Print list of reports and categories required for report data",
		Run: func(cmd *cobra.Command, args []string) {
			conf := &config.Config{}

			if credsFromEnv {
				readCredsFromEnv(conf)
			} else {
				askAuthQuestions(conf)
			}

			if err := fetchAndPrintReports(conf.ClientInfo(), filter); err != nil {
				log.Fatal(err)
			}

		},
	}

	cmd.Flags().BoolVar(&credsFromEnv, "creds-from-env", false, "Read credentials from envs instead of user input")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter reports containing the phrase")

	return cmd
}

func fetchAndPrintReports(info servicetitan.ClientInfo, filterTerm string) error {
	c, err := servicetitan.New(info)
	if err != nil {
		return err
	}

	log.Println("Fetching categories...")
	cat, err := c.ReportService.GetCategories(context.Background(), nil)
	if err != nil {
		return err
	}

	entries := []categoryReportEntry{}

	log.Println("Fetching reports...")
	for _, ctg := range cat.Items {
		rpt, err := fetchReportsForCategory(c, ctg)
		if err != nil {
			return err
		}

		entries = append(entries, categoryReportEntry{
			reports:  rpt,
			category: ctg,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Report ID", "Category ID", "Category Name", "Report Name"})

	for _, ent := range entries {
		for _, rpt := range ent.reports {
			if filterTerm != "" && !strings.Contains(rpt.Name, filterTerm) {
				continue
			}

			table.Append([]string{
				strconv.Itoa(rpt.ID),
				ent.category.ID,
				ent.category.Name,
				rpt.Name,
			})
		}
	}

	table.SetRowLine(true)
	table.Render()

	return nil
}

func fetchReportsForCategory(c *servicetitan.Client, category servicetitan.Category) ([]servicetitan.Report, error) {
	options := &servicetitan.ReportOptions{
		Page:     1,
		PageSize: 200,
	}

	reports := []servicetitan.Report{}

	for {
		col, err := c.ReportService.GetReports(context.Background(), category, options)
		if err != nil {
			return nil, err
		}
		reports = append(reports, col.Items...)

		if !col.HasMore {
			break
		}

		options.Page = col.Page + 1
	}

	return reports, nil
}
