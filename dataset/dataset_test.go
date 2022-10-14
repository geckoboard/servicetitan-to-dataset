package dataset

import (
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/servicetitan"
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
				"active":         {Type: "string", Optional: true, Name: "Active"},
				"completed_on":   {Type: "date", Optional: true, Name: "Completed date"},
				"name":           {Type: "string", Name: "Name"},
				"number_of_jobs": {Type: "number", Optional: true, Name: "Completed Jobs"},
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
				"active":         "TRUE",
				"completed_on":   "2021-10-13",
				"name":           "John Smith",
				"number_of_jobs": 5,
			},
			map[string]interface{}{
				"active":         "TRUE",
				"completed_on":   "2021-10-13",
				"name":           "Jane Doe",
				"number_of_jobs": 9,
			},
			map[string]interface{}{
				"active":         "FALSE",
				"completed_on":   "2021-10-13",
				"name":           "Hilary",
				"number_of_jobs": 15,
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
				[]interface{}{"John Smith", 5, true, "2021-10-13"},
				[]interface{}{"Jane Doe", 9, true, "2021-10-13"},
				[]interface{}{"Hilary", 15, false, "2021-10-13"},
			},
			Fields: []servicetitan.ReportField{
				{Name: "Name", Label: "Name", Type: "String"},
				{Name: "Number of jobs", Label: "Completed Jobs", Type: "Number"},
				{Name: "Active", Label: "Active", Type: "Boolean"},
				{Name: "Completed on", Label: "Completed date", Type: "Date"},
			},
			HasMore:  false,
			Page:     1,
			PageSize: 50,
		},
		DatasetOverrides: config.Dataset{
			RequiredFields: []string{"Name"},
		},
	}
}
