package dataset

import (
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/servicetitan"
	"strings"
	"testing"

	"github.com/jnormington/geckoboard"
	"gotest.tools/v3/assert"
)

func TestNewDatasetBuilder(t *testing.T) {
	report := &servicetitan.Report{}
	data := &servicetitan.ReportData{}

	t.Run("stores the config input correctly", func(t *testing.T) {
		got := NewDatasetBuilder(BuilderConfig{
			Report: report,
			Data:   data,
		})

		assert.Equal(t, got.data, data)
		assert.Equal(t, got.report, report)
	})
}

func TestDatasetBuilder_BuildSchema(t *testing.T) {
	t.Run("returns valid dataset schema", func(t *testing.T) {
		builder := NewDatasetBuilder(buildConfig())
		got := builder.BuildSchema()

		want := &geckoboard.Dataset{
			Name: "report_a",
			Fields: map[string]geckoboard.Field{
				"active":          {Type: "string", Optional: true, Name: "Active"},
				"completed_on":    {Type: "date", Optional: true, Name: "Completed date"},
				"name":            {Type: "string", Name: "Name"},
				"number_of_jobs":  {Type: "number", Optional: true, Name: "Completed Jobs"},
				"completion_rate": {Type: "percentage", Name: "Completion rate", Optional: true}, // type override in config
				"created_on":      {Type: "datetime", Optional: true, Name: "Created on"},
				"tags":            {Type: "string", Optional: true, Name: "Tags"},
			},
			UniqueBy: []string{"name"},
		}

		assert.DeepEqual(t, got, want)
	})

	t.Run("removes invalid characters from the report name for the dataset name", func(t *testing.T) {
		conf := buildConfig()
		conf.Report.Name = "My report i$ the best 1235"

		builder := NewDatasetBuilder(conf)
		got := builder.BuildSchema()
		assert.Equal(t, got.Name, "my_report_i_the_best_1235")
	})

	t.Run("returns report name as the dataset name when override is empty", func(t *testing.T) {
		builder := NewDatasetBuilder(buildConfig())
		got := builder.BuildSchema()
		assert.Equal(t, got.Name, "report_a")
	})

	t.Run("returns valid dataset name from the overridden value", func(t *testing.T) {
		specs := []struct {
			in  string
			out string
		}{
			{"dataset-2", "dataset-2"},
			{"dataset abc.2 - v3", "dataset_abc.2_-_v3"},
		}

		for _, tc := range specs {
			t.Run(tc.in, func(t *testing.T) {
				conf := buildConfig()
				conf.DatasetOverrides.Name = tc.in

				builder := NewDatasetBuilder(conf)
				got := builder.BuildSchema()
				assert.Equal(t, got.Name, tc.out)
			})
		}
	})
}

func TestDatasetBuilder_BuildData(t *testing.T) {
	t.Run("returns the dataset in the correct format", func(t *testing.T) {
		builder := NewDatasetBuilder(buildConfig())
		got := builder.BuildData()

		assert.DeepEqual(t, got, geckoboard.Data{
			map[string]interface{}{
				"active":          "TRUE",
				"completed_on":    "2021-10-13",
				"name":            "John Smith",
				"number_of_jobs":  5,
				"completion_rate": 0.12,
				"created_on":      "2023-10-13T00:00:00-05:00",
				"tags":            strings.Repeat("a", 256),
			},
			map[string]interface{}{
				"active":          "TRUE",
				"completed_on":    "2021-10-13",
				"name":            "Jane Doe",
				"number_of_jobs":  9,
				"completion_rate": 0.24,
				"created_on":      "2023-10-13T00:00:00-05:00Z",
				"tags":            strings.Repeat("b", 100),
			},
			map[string]interface{}{
				"active":          "FALSE",
				"completed_on":    "2021-10-13",
				"name":            "Hilary",
				"number_of_jobs":  15,
				"completion_rate": 0.87,
				"created_on":      "2023-10-13T00:00:00-05:00",
				"tags":            strings.Repeat("界", 256),
			},
		})
	})
}

func buildConfig() BuilderConfig {
	return BuilderConfig{
		Report: &servicetitan.Report{
			ID:   2222222,
			Name: "Report A",
			Fields: []servicetitan.ReportField{
				{Name: "Name", Label: "Name", Type: "String"},
				{Name: "Number of jobs", Label: "Completed Jobs", Type: "Number"},
				{Name: "Active", Label: "Active", Type: "Boolean"},
				{Name: "Completed on", Label: "Completed date", Type: "Date"},
				{Name: "Completion rate", Label: "Completion rate", Type: "Number"},
				{Name: "Created on", Label: "Created on", Type: "Datetime"},
				{Name: "Tags", Label: "Tags", Type: "String"},
			},
			Parameters: []servicetitan.ReportParameter{
				{Name: "From", Label: "From", DataType: "Date", IsRequired: true},
				{Name: "To", Label: "To", DataType: "Date", IsRequired: true},
				{
					Name:     "IncludeInactive",
					Label:    "Include Inactive Technicians",
					DataType: "Boolean",
				},
			},
		},
		Data: &servicetitan.ReportData{
			Data: []interface{}{
				[]interface{}{"John Smith", 5, true, "2021-10-13", 0.12, "2023-10-13T00:00:00-05:00", strings.Repeat("a", 300)},
				[]interface{}{"Jane Doe", 9, true, "2021-10-13", 0.24, "2023-10-13T00:00:00-05:00Z", strings.Repeat("b", 100)},
				[]interface{}{"Hilary", 15, false, "2021-10-13", 0.87, "2023-10-13T00:00:00-05:00", strings.Repeat("界", 300)},
			},
			Fields: []servicetitan.ReportField{
				{Name: "Name", Label: "Name", Type: "String"},
				{Name: "Number of jobs", Label: "Completed Jobs", Type: "Number"},
				{Name: "Active", Label: "Active", Type: "Boolean"},
				{Name: "Completed on", Label: "Completed date", Type: "Date"},
				{Name: "Completion rate", Label: "Completion rate", Type: "Number"},
				{Name: "Created on", Label: "Created on", Type: "Datetime"},
				{Name: "Tags", Label: "Tags", Type: "String"},
			},
			HasMore:  false,
			Page:     1,
			PageSize: 50,
		},
		DatasetOverrides: config.Dataset{
			RequiredFields: []string{"Name"},
			FieldOverrides: []config.ReportField{
				{
					Name: "Completion rate",
					Type: "Percentage",
				},
				{
					Name: "Created on",
					Type: "Datetime",
				},
			},
		},
	}
}
