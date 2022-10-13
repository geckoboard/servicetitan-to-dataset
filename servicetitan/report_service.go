package servicetitan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type ReportService interface {
	GetCategories(context.Context, *PaginationOptions) (*CategoryList, error)
	GetReports(context.Context, Category, *PaginationOptions) (*ReportList, error)
	GetReport(_ context.Context, categoryID, reportID string) (*Report, error)
	GetReportData(context.Context, ReportDataRequest, *PaginationOptions) (*ReportData, error)
}

type reportService struct {
	baseURL string
	client  *Client
}

func (r reportService) GetCategories(ctx context.Context, options *PaginationOptions) (*CategoryList, error) {
	url := r.client.buildURL(r.baseURL, "/report-categories", r.buildParams(options))
	req, err := r.client.buildGETRequest(url)
	if err != nil {
		return nil, err
	}

	list := &CategoryList{}
	if err := r.client.doRequest(req.WithContext(ctx), list); err != nil {
		return nil, err
	}

	return list, nil
}

func (r reportService) GetReports(ctx context.Context, category Category, options *PaginationOptions) (*ReportList, error) {
	path := fmt.Sprintf("/report-category/%s/reports", category.ID)
	url := r.client.buildURL(r.baseURL, path, r.buildParams(options))
	req, err := r.client.buildGETRequest(url)
	if err != nil {
		return nil, err
	}

	list := &ReportList{}
	if err := r.client.doRequest(req.WithContext(ctx), list); err != nil {
		return nil, err
	}

	return list, nil
}

func (r reportService) GetReport(ctx context.Context, categoryID, reportID string) (*Report, error) {
	path := fmt.Sprintf("/report-category/%s/reports/%s", categoryID, reportID)
	url := r.client.buildURL(r.baseURL, path, nil)

	req, err := r.client.buildGETRequest(url)
	if err != nil {
		return nil, err
	}

	report := &Report{}
	if err := r.client.doRequest(req.WithContext(ctx), report); err != nil {
		return nil, err
	}

	return report, nil
}

func (r reportService) GetReportData(ctx context.Context, opts ReportDataRequest, pagination *PaginationOptions) (*ReportData, error) {
	path := fmt.Sprintf("/report-category/%s/reports/%s/data", opts.CategoryID, opts.ReportID)
	url := r.client.buildURL(r.baseURL, path, r.buildParams(pagination))

	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	req, err := r.client.buildPOSTRequest(url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	data := &ReportData{}
	if err := r.client.doRequest(req.WithContext(ctx), data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r reportService) buildParams(options *PaginationOptions) url.Values {
	params := url.Values{"page": {"1"}, "pageSize": {"50"}}

	if options != nil {
		params = url.Values{
			"page":     {strconv.Itoa(options.Page)},
			"pageSize": {strconv.Itoa(options.PageSize)},
		}
	}

	return params
}
