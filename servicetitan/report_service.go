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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ReportList struct {
	Items []Report `json:"data"`

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
