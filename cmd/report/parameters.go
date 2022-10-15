package report

import (
	"context"
	"fmt"
	"log"
	"os"
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/servicetitan"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func ParametersCommand() *cobra.Command {
	var (
		reportID   string
		categoryID string
	)

	cmd := &cobra.Command{
		Use:   "parameters",
		Short: "Print parameters for a report",
		Run: func(cmd *cobra.Command, args []string) {
			if reportID == "" || categoryID == "" {
				log.Fatal("both --report and --category are required, you can get from using 'reports list' command")
			}

			cfg, err := config.LoadFile(cmd.Flag("config").Value.String())
			if err != nil {
				log.Fatal(err)
			}

			if err := fetchAndDisplayParameters(cfg.ServiceTitan, categoryID, reportID); err != nil {
				log.Fatal(err)
			}

		},
	}

	cmd.Flags().StringVar(&reportID, "report", "", "Report ID to fetch the report for parameters")
	cmd.Flags().StringVar(&categoryID, "category", "", "Category ID to fetch the report parameters")

	return cmd
}

func fetchAndDisplayParameters(cfg config.ServiceTitan, categoryID, reportID string) error {
	c, err := servicetitan.New(cfg)
	if err != nil {
		return err
	}

	report, err := c.ReportService.GetReport(context.Background(), categoryID, reportID)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("Report id: ", report.ID)
	fmt.Println("Report name: ", report.Name)

	paramTable := tablewriter.NewWriter(os.Stdout)
	paramTable.SetRowLine(true)
	paramTable.SetHeader([]string{"Paramter name", "Label", "Data type", "Array?", "Required?"})

	fieldTable := tablewriter.NewWriter(os.Stdout)
	fieldTable.SetRowLine(true)
	fieldTable.SetHeader([]string{"Field Name", "Label", "Type"})

	for _, param := range report.Parameters {
		paramTable.Append([]string{
			param.Name,
			param.Label,
			param.DataType,
			strings.ToUpper(strconv.FormatBool(param.IsArray)),
			strings.ToUpper(strconv.FormatBool(param.IsRequired)),
		})
	}

	for _, field := range report.Fields {
		fieldTable.Append([]string{field.Label, field.Name, field.Type})
	}

	fmt.Println("")
	fmt.Println("Report fields:")
	fieldTable.Render()

	fmt.Println("")
	fmt.Println("Report parameters:")
	paramTable.Render()

	return nil
}
