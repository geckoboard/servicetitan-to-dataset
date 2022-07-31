package cmd

import (
	"context"
	"log"
	"os"
	"servicetitan-to-dataset/servicetitan"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type categoryReportEntry struct {
	reports  []servicetitan.Report
	category servicetitan.Category
}

func ListReportsCommand() *cobra.Command {
	var credsFromEnv bool

	cmd := &cobra.Command{
		Use:   "list-reports",
		Short: "Print list of reports and categories required for report data",
		Run: func(cmd *cobra.Command, args []string) {
			conf := &servicetitan.ClientInfo{}

			if credsFromEnv {
				readCredsFromEnv(conf)
			} else {
				askAuthQuestions(conf)
			}

			if err := fetchAndPrintReports(*conf); err != nil {
				log.Fatal(err)
			}

		},
	}

	cmd.Flags().BoolVar(&credsFromEnv, "creds-from-env", false, "Read credentials from envs instead of user input")

	return cmd
}

func fetchAndPrintReports(info servicetitan.ClientInfo) error {
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
	table.SetHeader([]string{"Category ID", "Category Name", "Report ID", "Report Name"})

	for _, ent := range entries {
		for _, rpt := range ent.reports {
			table.Append([]string{
				ent.category.ID,
				ent.category.Name,
				strconv.Itoa(rpt.ID),
				rpt.Name,
			})
		}
	}

	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
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
