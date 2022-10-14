package processor

import (
	"context"
	"fmt"
	"regexp"
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/dataset"
	"servicetitan-to-dataset/servicetitan"
	"strconv"
	"strings"
	"time"

	"github.com/jnormington/geckoboard"
)

const dateFormat = "2006-01-02"

var nowSubRegexp = regexp.MustCompile(`NOW(\-|\+)(\d+)`)

type ReportProcessor struct {
	maxDatasetRecords int
	config            *config.Config
	timeNow           func() time.Time

	serviceTitanClient *servicetitan.Client
	geckoboardClient   *geckoboard.Client
}

func New(cfg *config.Config) ReportProcessor {
	c, _ := servicetitan.New(cfg.ServiceTitan)
	gb := geckoboard.New("https://api.geckoboard.com", cfg.Geckoboard.APIKey)

	return ReportProcessor{
		maxDatasetRecords:  5000,
		config:             cfg,
		serviceTitanClient: c,
		geckoboardClient:   gb,
		timeNow:            time.Now,
	}
}

func (r ReportProcessor) Process(ctx context.Context, entry config.Entry) error {
	report, err := r.serviceTitanClient.ReportService.GetReport(ctx, entry.Report.CategoryID, entry.Report.ID)
	if err != nil {
		return err
	}

	data, err := r.fetchReportData(ctx, report, entry)
	if err != nil {
		return err
	}

	builder := dataset.NewDatasetBuilder(dataset.BuilderConfig{
		Report:           report,
		Data:             data,
		DatasetOverrides: entry.Dataset,
	})

	schema := builder.BuildSchema()
	if err := r.geckoboardClient.DatasetService.FindOrCreate(ctx, schema); err != nil {
		return err
	}

	if strings.ToLower(entry.Dataset.Type) == "append" {
		return r.geckoboardClient.DatasetService.AppendData(ctx, schema, builder.BuildData())
	}

	return r.geckoboardClient.DatasetService.ReplaceData(ctx, schema, builder.BuildData())
}

func (r *ReportProcessor) fetchReportData(ctx context.Context, report *servicetitan.Report, entry config.Entry) (*servicetitan.ReportData, error) {
	reportParams, err := r.buildReportParameters(report, entry)
	if err != nil {
		return nil, err
	}

	reportOpts := servicetitan.ReportDataRequest{
		CategoryID: entry.Report.CategoryID,
		ReportID:   entry.Report.ID,
		Parameters: reportParams,
	}

	reportData := &servicetitan.ReportData{}

	// We can only fetch a single page - with a rate limit of 1 per 5 minutes
	pagination := &servicetitan.PaginationOptions{Page: 1, PageSize: 5000}
	resp, err := r.serviceTitanClient.ReportService.GetReportData(ctx, reportOpts, pagination)
	if err != nil {
		return nil, err
	}

	reportData.Data = append(reportData.Data, resp.Data...)
	if reportData.Fields == nil {
		reportData.Fields = resp.Fields
	}

	return reportData, nil
}

// Build the parameters from the config to servicetitan compatible parameters.
// This also supports special NOW and NOW-n keywords for date fields - which if
// the fields are of type Date will be replaced with the current time and current time -n
func (r *ReportProcessor) buildReportParameters(report *servicetitan.Report, ent config.Entry) ([]servicetitan.DataRequestParamters, error) {
	params := []servicetitan.DataRequestParamters{}

	for _, p := range ent.Report.Parameters {
		param := r.lookupParameter(report, p.Name)

		if param == nil {
			return nil, fmt.Errorf("invalid param %q for report %v", p.Name, report.ID)
		}

		value := p.Value

		if param.DataType == "Date" && p.Value == "NOW" {
			value = r.timeNow().Format(dateFormat)
		}

		val, _ := value.(string)
		if param.DataType == "Date" && nowSubRegexp.MatchString(val) {
			matches := nowSubRegexp.FindStringSubmatch(val)

			if len(matches) == 3 {
				parsedNum, _ := strconv.Atoi(matches[2])
				days := time.Duration(parsedNum*24) * time.Hour

				switch matches[1] {
				case "-":
					value = r.timeNow().Add(-days).Format(dateFormat)
				case "+":
					value = r.timeNow().Add(days).Format(dateFormat)
				}
			}

		}

		params = append(params, servicetitan.DataRequestParamters{Name: p.Name, Value: value})
	}

	return params, nil
}

func (r *ReportProcessor) lookupParameter(report *servicetitan.Report, key string) *servicetitan.ReportParameter {
	for _, p := range report.Parameters {
		if p.Name == key {
			return &p
		}
	}

	return nil
}
