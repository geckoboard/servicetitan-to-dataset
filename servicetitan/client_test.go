package servicetitan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"servicetitan-to-dataset/config"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestClient_New(t *testing.T) {
	t.Run("returns new client", func(t *testing.T) {
		cfg := config.ServiceTitan{TenantID: "tenant_123"}
		c, err := New(cfg)
		assert.NilError(t, err)

		assert.Assert(t, c.client != nil)
		assert.Assert(t, c.session == nil)
		assert.Equal(t, c.config, cfg)

		authSrv := c.AuthService.(authService)
		assert.Equal(t, authSrv.client, c)
		assert.Equal(t, authSrv.baseURL, "https://auth.servicetitan.io")

		reportSrv := c.ReportService.(reportService)
		assert.Equal(t, reportSrv.client, c)
		assert.Equal(t, reportSrv.baseURL, "https://api.servicetitan.io/reporting/v2/tenant/tenant_123")
	})
}

func TestClient_Authorization(t *testing.T) {
	t.Run("creates new session on first request", func(t *testing.T) {
		authCalls := 0
		reportCalls := 0

		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			reportCalls++
			assert.Equal(t, r.Header.Get("Authorization"), "tok_1234")
			io.WriteString(w, "{}")
		})

		c := &Client{client: http.DefaultClient}
		c.AuthService = &mockAuthService{
			getTokenFn: func() (*Session, error) {
				authCalls++
				return &Session{
					Token:     "tok_1234",
					ExpiresAt: time.Now().UTC(),
				}, nil
			},
		}
		c.ReportService = mockReportService{
			client:  c,
			baseURL: server.URL,
		}

		_, err := c.ReportService.GetCategories(nil, nil)
		assert.NilError(t, err)

		assert.Equal(t, authCalls, 1)
		assert.Equal(t, reportCalls, 1)
	})

	t.Run("uses same session on subsequent requests", func(t *testing.T) {
		authCalls := 0
		reportCalls := 0

		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			reportCalls++
			assert.Equal(t, r.Header.Get("Authorization"), "tok_1234")
			io.WriteString(w, "{}")
		})

		c := &Client{client: http.DefaultClient}
		c.AuthService = &mockAuthService{
			getTokenFn: func() (*Session, error) {
				authCalls++
				return &Session{
					Token:     "tok_1234",
					ExpiresAt: time.Now().UTC().Add(2 * time.Minute),
				}, nil
			},
		}
		c.ReportService = mockReportService{
			client:  c,
			baseURL: server.URL,
		}

		_, err := c.ReportService.GetCategories(nil, nil)
		assert.NilError(t, err)
		_, err = c.ReportService.GetCategories(nil, nil)
		assert.NilError(t, err)
		_, err = c.ReportService.GetReports(nil, Category{}, nil)
		assert.NilError(t, err)

		assert.Equal(t, authCalls, 1)
		assert.Equal(t, reportCalls, 3)
	})

	t.Run("refreshes auth token when expired", func(t *testing.T) {
		authCalls := 0
		reportCalls := 0

		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			reportCalls++

			if reportCalls > 2 {
				assert.Equal(t, r.Header.Get("Authorization"), "tok_1232")
			} else {
				assert.Equal(t, r.Header.Get("Authorization"), "tok_1231")
			}
			io.WriteString(w, "{}")
		})

		c := &Client{client: http.DefaultClient}
		c.AuthService = &mockAuthService{
			getTokenFn: func() (*Session, error) {
				authCalls++

				return &Session{
					Token:     fmt.Sprintf("tok_123%d", authCalls),
					ExpiresAt: time.Now().UTC().Add(2 * time.Minute),
				}, nil
			},
		}
		c.ReportService = mockReportService{
			client:  c,
			baseURL: server.URL,
		}

		_, err := c.ReportService.GetCategories(nil, nil)
		assert.NilError(t, err)
		_, err = c.ReportService.GetReports(nil, Category{}, nil)
		assert.NilError(t, err)

		assert.Equal(t, authCalls, 1)
		assert.Equal(t, reportCalls, 2)

		c.session.ExpiresAt = time.Now().UTC()
		_, err = c.ReportService.GetCategories(nil, nil)
		assert.NilError(t, err)

		assert.Equal(t, authCalls, 2)
		assert.Equal(t, reportCalls, 3)
	})
}

type mockAuthService struct {
	getTokenFn func() (*Session, error)
	calls      int
}

func (m *mockAuthService) GetToken(context.Context, config.ServiceTitan) (*Session, error) {
	m.calls++

	if m.getTokenFn != nil {
		return m.getTokenFn()
	}

	return &Session{
		Token:     fmt.Sprintf("tok_123%d", m.calls),
		ExpiresAt: time.Now().UTC().Add(1 * time.Minute),
	}, nil
}

type mockReportService struct {
	client  *Client
	baseURL string
}

func (r mockReportService) GetCategories(context.Context, *PaginationOptions) (*CategoryList, error) {
	url, _ := url.Parse(r.baseURL)
	req := &http.Request{URL: url, Header: http.Header{}}
	r.client.doRequest(req, nil)
	return nil, nil
}

func (r mockReportService) GetReports(context.Context, Category, *PaginationOptions) (*ReportList, error) {
	url, _ := url.Parse(r.baseURL)
	req := &http.Request{URL: url, Header: http.Header{}}
	r.client.doRequest(req, nil)
	return nil, nil
}

func (r mockReportService) GetReport(context.Context, string, string) (*Report, error) {
	url, _ := url.Parse(r.baseURL)
	req := &http.Request{URL: url, Header: http.Header{}}
	r.client.doRequest(req, nil)
	return nil, nil
}

func (r mockReportService) GetReportData(context.Context, ReportDataRequest, *PaginationOptions) (*ReportData, error) {
	url, _ := url.Parse(r.baseURL)
	req := &http.Request{URL: url, Header: http.Header{}}
	r.client.doRequest(req, nil)
	return nil, nil
}
