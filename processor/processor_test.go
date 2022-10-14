package processor

import (
	"context"
	"errors"
	"servicetitan-to-dataset/config"
	"servicetitan-to-dataset/servicetitan"
	"testing"
	"time"

	"github.com/jnormington/geckoboard"
	"gotest.tools/v3/assert"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.Config{
		ServiceTitan: config.ServiceTitan{},
		Geckoboard:   config.Geckoboard{},
	}

	out := New(cfg)

	assert.Equal(t, out.maxDatasetRecords, 5000)
	assert.Equal(t, out.config, cfg)
	assert.Assert(t, out.serviceTitanClient != nil)
	assert.Assert(t, out.geckoboardClient != nil)
}

func TestProcessor_Processor(t *testing.T) {
	t.Run("queries the correct report and report data", func(t *testing.T) {
		var (
			calledReport     bool
			calledReportData bool
		)
		proc, rs, _ := buildProcessorWithMocks()

		rs.getReportFn = func(gotCategoryID, gotReportID string) (*servicetitan.Report, error) {
			assert.Equal(t, gotCategoryID, "category-abc")
			assert.Equal(t, gotReportID, "1234")
			calledReport = true

			return &servicetitan.Report{
				ID:   1234,
				Name: "Report A",
				Fields: []servicetitan.ReportField{
					{Name: "Name", Label: "Name", Type: "String"},
					{Name: "Number of jobs", Label: "Completed Jobs", Type: "Number"},
					{Name: "Active", Label: "Active", Type: "Boolean"},
					{Name: "Completed on", Label: "Completed date", Type: "Date"},
				},
			}, nil
		}

		rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
			assert.DeepEqual(t, gotPagination, &servicetitan.PaginationOptions{
				PageSize: 5000,
				Page:     1,
			})

			assert.DeepEqual(t, got, servicetitan.ReportDataRequest{
				ReportID:   "1234",
				CategoryID: "category-abc",
				Parameters: []servicetitan.DataRequestParamters{},
			})

			calledReportData = true
			return &servicetitan.ReportData{
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
			}, nil
		}

		err := proc.Process(context.Background(), config.Entry{
			Report: config.Report{
				ID:         "1234",
				CategoryID: "category-abc",
			},
		})

		assert.NilError(t, err)
		assert.Assert(t, calledReport)
		assert.Assert(t, calledReportData)
	})

	t.Run("queries report data with user set parameters", func(t *testing.T) {
		var calledReportData bool
		proc, rs, _ := buildProcessorWithMocks()

		rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
			assert.DeepEqual(t, got, servicetitan.ReportDataRequest{
				ReportID:   "1234",
				CategoryID: "category-abc",
				Parameters: []servicetitan.DataRequestParamters{
					{
						Name:  "From",
						Value: "2021-10-13",
					},
					{
						Name:  "To",
						Value: "2021-10-14",
					},
				},
			})

			calledReportData = true
			return &servicetitan.ReportData{
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
			}, nil
		}

		err := proc.Process(context.Background(), config.Entry{
			Report: config.Report{
				ID:         "1234",
				CategoryID: "category-abc",
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
		})

		assert.NilError(t, err)
		assert.Assert(t, calledReportData)
	})

	t.Run("pushes the correct schema/data to geckoboard", func(t *testing.T) {
		proc, _, ds := buildProcessorWithMocks()

		var (
			calledFindOrCreate bool
			calledReplaceData  bool
		)

		ds.findOrCreateFn = func(got *geckoboard.Dataset) error {
			assert.DeepEqual(t, got, &geckoboard.Dataset{
				Name: "report_a",
				Fields: map[string]geckoboard.Field{
					"name":           {Type: "string", Name: "Name"},
					"completed_on":   {Type: "date", Name: "Completed date", Optional: true},
					"number_of_jobs": {Type: "number", Name: "Completed Jobs", Optional: true},
					"active":         {Type: "string", Name: "Active", Optional: true},
				},
				UniqueBy: []string{"name"},
			})

			calledFindOrCreate = true
			return nil
		}

		ds.replaceDataFn = func(ds *geckoboard.Dataset, got geckoboard.Data) error {
			assert.DeepEqual(t, got, geckoboard.Data{
				{
					"active":         "TRUE",
					"completed_on":   "2021-10-13",
					"name":           "John Smith",
					"number_of_jobs": 5,
				},
				{
					"active":         "TRUE",
					"completed_on":   "2021-10-13",
					"name":           "Jane Doe",
					"number_of_jobs": 9,
				},
				{
					"active":         "FALSE",
					"completed_on":   "2021-10-13",
					"name":           "Hilary",
					"number_of_jobs": 15,
				},
			})

			calledReplaceData = true
			return nil
		}

		err := proc.Process(context.Background(), config.Entry{
			Dataset: config.Dataset{
				RequiredFields: []string{"Name"},
			},
		})
		assert.NilError(t, err)
		assert.Assert(t, calledFindOrCreate)
		assert.Assert(t, calledReplaceData)
	})

	t.Run("dynamic date parameter value", func(t *testing.T) {
		t.Run("replaces NOW with the current date and NOW-1 with yesterday", func(t *testing.T) {
			var calledReportData bool
			proc, rs, _ := buildProcessorWithMocks()

			rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
				assert.DeepEqual(t, got.Parameters, []servicetitan.DataRequestParamters{
					{Name: "From", Value: "2022-06-06"},
					{Name: "To", Value: "2022-06-07"},
				})
				calledReportData = true

				return &servicetitan.ReportData{
					Data:   []interface{}{},
					Fields: []servicetitan.ReportField{},
				}, nil
			}

			err := proc.Process(context.Background(), config.Entry{
				Report: config.Report{
					ID:         "1234",
					CategoryID: "category-abc",
					Parameters: []config.Parameter{
						{
							Name:  "From",
							Value: "NOW-1",
						},
						{
							Name:  "To",
							Value: "NOW",
						},
					},
				},
			})

			assert.NilError(t, err)
			assert.Assert(t, calledReportData)
		})

		t.Run("doesnt replace NOW when the parameter type is not a date", func(t *testing.T) {
			var calledReportData bool
			proc, rs, _ := buildProcessorWithMocks()

			rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
				assert.DeepEqual(t, got.Parameters, []servicetitan.DataRequestParamters{{Name: "Username", Value: "NOW"}})
				calledReportData = true

				return &servicetitan.ReportData{
					Data:   []interface{}{},
					Fields: []servicetitan.ReportField{},
				}, nil
			}

			err := proc.Process(context.Background(), config.Entry{
				Report: config.Report{
					Parameters: []config.Parameter{
						{
							Name:  "Username",
							Value: "NOW",
						},
					},
				},
			})

			assert.NilError(t, err)
			assert.Assert(t, calledReportData)
		})

		t.Run("doesnt replace NOW- when incomplete", func(t *testing.T) {
			var calledReportData bool
			proc, rs, _ := buildProcessorWithMocks()

			rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
				assert.DeepEqual(t, got.Parameters, []servicetitan.DataRequestParamters{{Name: "From", Value: "NOW-"}})
				calledReportData = true

				return &servicetitan.ReportData{
					Data:   []interface{}{},
					Fields: []servicetitan.ReportField{},
				}, nil
			}

			err := proc.Process(context.Background(), config.Entry{
				Report: config.Report{
					Parameters: []config.Parameter{
						{
							Name:  "From",
							Value: "NOW-",
						},
					},
				},
			})

			assert.NilError(t, err)
			assert.Assert(t, calledReportData)
		})

		t.Run("replaces NOW+1 when today +1 day", func(t *testing.T) {
			var calledReportData bool
			proc, rs, _ := buildProcessorWithMocks()

			rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
				assert.DeepEqual(t, got.Parameters, []servicetitan.DataRequestParamters{
					{Name: "From", Value: "2022-06-08"},
					{Name: "To", Value: "2022-06-09"},
				})
				calledReportData = true

				return &servicetitan.ReportData{
					Data:   []interface{}{},
					Fields: []servicetitan.ReportField{},
				}, nil
			}

			err := proc.Process(context.Background(), config.Entry{
				Report: config.Report{
					ID:         "1234",
					CategoryID: "category-abc",
					Parameters: []config.Parameter{
						{
							Name:  "From",
							Value: "NOW+1",
						},
						{
							Name:  "To",
							Value: "NOW+2",
						},
					},
				},
			})

			assert.NilError(t, err)
			assert.Assert(t, calledReportData)
		})
	})

	t.Run("returns error when report fetch fails", func(t *testing.T) {
		proc, rs, _ := buildProcessorWithMocks()

		rs.getReportFn = func(gotCategoryID, gotReportID string) (*servicetitan.Report, error) {
			return nil, errors.New("report fetch failed")
		}

		err := proc.Process(context.Background(), config.Entry{})
		assert.ErrorContains(t, err, "report fetch failed")
	})

	t.Run("returns error when fetching report data fails", func(t *testing.T) {
		proc, rs, _ := buildProcessorWithMocks()

		rs.getReportDataFn = func(got servicetitan.ReportDataRequest, gotPagination *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
			return nil, errors.New("missing parameters")
		}

		err := proc.Process(context.Background(), config.Entry{})
		assert.ErrorContains(t, err, "missing parameters")
	})

	t.Run("returns error when find or create dataset fails", func(t *testing.T) {
		proc, _, ds := buildProcessorWithMocks()

		ds.findOrCreateFn = func(*geckoboard.Dataset) error {
			return errors.New("fetch dataset error")
		}

		err := proc.Process(context.Background(), config.Entry{})
		assert.ErrorContains(t, err, "fetch dataset error")
	})

	t.Run("returns error when replacing data call fails", func(t *testing.T) {
		proc, _, ds := buildProcessorWithMocks()

		ds.replaceDataFn = func(*geckoboard.Dataset, geckoboard.Data) error {
			return errors.New("replace data error")
		}

		err := proc.Process(context.Background(), config.Entry{})
		assert.ErrorContains(t, err, "replace data error")
	})

	t.Run("returns error when appending data call fails", func(t *testing.T) {
		proc, _, ds := buildProcessorWithMocks()

		ds.appendDataFn = func(*geckoboard.Dataset, geckoboard.Data) error {
			return errors.New("append data error")
		}

		err := proc.Process(context.Background(), config.Entry{
			Dataset: config.Dataset{Type: "append"},
		})
		assert.ErrorContains(t, err, "append data error")
	})

	t.Run("returns error when invalid param in config", func(t *testing.T) {
		proc, _, _ := buildProcessorWithMocks()

		err := proc.Process(context.Background(), config.Entry{
			Report: config.Report{
				Parameters: []config.Parameter{
					{
						Name:  "not-existing",
						Value: "blah",
					},
				},
			},
			Dataset: config.Dataset{Type: "append"},
		})
		assert.ErrorContains(t, err, "invalid param \"not-existing\" for report 2222222")
	})
}

func buildProcessorWithMocks() (ReportProcessor, *mockReportService, *mockDatasetService) {
	reportSrv := &mockReportService{}
	datasetSrv := &mockDatasetService{}

	proc := New(&config.Config{
		ServiceTitan: config.ServiceTitan{},
		Geckoboard:   config.Geckoboard{},
	})

	proc.timeNow = func() time.Time {
		return time.Date(2022, 6, 7, 8, 11, 0, 0, time.UTC)
	}
	proc.serviceTitanClient.ReportService = reportSrv
	proc.geckoboardClient.DatasetService = datasetSrv

	return proc, reportSrv, datasetSrv
}

type mockDatasetService struct {
	findOrCreateFn func(*geckoboard.Dataset) error
	appendDataFn   func(*geckoboard.Dataset, geckoboard.Data) error
	replaceDataFn  func(*geckoboard.Dataset, geckoboard.Data) error
}

func (m *mockDatasetService) FindOrCreate(_ context.Context, ds *geckoboard.Dataset) error {
	if m.findOrCreateFn == nil {
		return nil
	}

	return m.findOrCreateFn(ds)
}

func (m *mockDatasetService) AppendData(_ context.Context, ds *geckoboard.Dataset, d geckoboard.Data) error {
	if m.appendDataFn == nil {
		return nil
	}
	return m.appendDataFn(ds, d)
}

func (m *mockDatasetService) ReplaceData(_ context.Context, ds *geckoboard.Dataset, d geckoboard.Data) error {
	if m.replaceDataFn == nil {
		return nil
	}

	return m.replaceDataFn(ds, d)
}

type mockReportService struct {
	getReportFn     func(string, string) (*servicetitan.Report, error)
	getReportDataFn func(servicetitan.ReportDataRequest, *servicetitan.PaginationOptions) (*servicetitan.ReportData, error)
}

func (r *mockReportService) GetCategories(context.Context, *servicetitan.PaginationOptions) (*servicetitan.CategoryList, error) {
	return nil, errors.New("not expected to be called")
}

func (r *mockReportService) GetReports(context.Context, servicetitan.Category, *servicetitan.PaginationOptions) (*servicetitan.ReportList, error) {
	return nil, errors.New("not expected to be called")
}

func (r *mockReportService) GetReport(_ context.Context, categoryID, reportID string) (*servicetitan.Report, error) {
	if r.getReportFn == nil {
		return &servicetitan.Report{
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
				{Name: "Username", Label: "User", DataType: "String", IsRequired: false},
				{
					Name:     "IncludeInactive",
					Label:    "Include Inactive Technicians",
					DataType: "Boolean",
				},
			},
		}, nil
	}

	return r.getReportFn(categoryID, reportID)
}

func (r *mockReportService) GetReportData(_ context.Context, rdr servicetitan.ReportDataRequest, po *servicetitan.PaginationOptions) (*servicetitan.ReportData, error) {
	if r.getReportDataFn == nil {
		return &servicetitan.ReportData{
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
		}, nil
	}

	return r.getReportDataFn(rdr, po)
}
