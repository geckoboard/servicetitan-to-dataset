package report

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

func ListCommand() *cobra.Command {
	var reportsFilter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists reports across all categories",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadFile(cmd.Flag("config").Value.String())
			if err != nil {
				log.Fatal(err)
			}

			if err := fetchAndPrintReports(cfg.ServiceTitan, reportsFilter); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&reportsFilter, "filter", "", "Filter list of reports containing the specific phrase")

	return cmd
}

func fetchAndPrintReports(cfg config.ServiceTitan, filterTerm string) error {
	c, err := servicetitan.New(cfg)
	if err != nil {
		return err
	}

	log.Println("Fetching categories...")
	cat, err := c.ReportService.GetCategories(context.Background(), nil)
	if err != nil {
		return err
	}

	entries := []categoryReportEntry{}

	for _, ctg := range cat.Items {
		log.Println("Fetching reports for category", ctg.Name, "...")
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
