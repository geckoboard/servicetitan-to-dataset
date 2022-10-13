package servicetitan

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type ReportService interface {
	GetCategories(context.Context, *ReportOptions) (*CategoryList, error)
	GetReports(context.Context, Category, *ReportOptions) (*ReportList, error)
	GetReport(_ context.Context, categoryID, reportID string) (*Report, error)
}

type reportService struct {
	baseURL string
	client  *Client
}

type ReportOptions struct {
	PageSize int
	Page     int
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryList struct {
	Items []Category `json:"data"`

	HasMore  bool `json:"hasMore"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"totalCount"`
}

type Report struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Fields     []ReportField     `json:"fields,omitempty"`
	Parameters []ReportParameter `json:"parameters,omitempty"`
}

type ReportList struct {
	Items []Report `json:"data"`

	HasMore  bool `json:"hasMore"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"totalCount"`
}

type ReportField struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"dataType"`
}

type ReportParameter struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	DataType   string `json:"dataType"`
	IsArray    bool   `json:"isArray"`
	IsRequired bool   `json:"isRequired"`
}

type ReportData struct {
	Data   []interface{} `json:"data"`
	Fields []ReportField `json:"fields"`

	HasMore  bool `json:"hasMore"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Total    int  `json:"totalCount"`
}

func (r reportService) GetCategories(ctx context.Context, options *ReportOptions) (*CategoryList, error) {
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

func (r reportService) GetReports(ctx context.Context, category Category, options *ReportOptions) (*ReportList, error) {
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

func (r reportService) buildParams(options *ReportOptions) url.Values {
	params := url.Values{"page": {"1"}, "pageSize": {"50"}}

	if options != nil {
		params = url.Values{
			"page":     {strconv.Itoa(options.Page)},
			"pageSize": {strconv.Itoa(options.PageSize)},
		}
	}

	return params
}
