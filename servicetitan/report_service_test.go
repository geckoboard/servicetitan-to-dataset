package servicetitan

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestReportService_GetCategories(t *testing.T) {
	t.Run("fetches first page of categories when report options is nil", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			assert.NilError(t, r.ParseForm())
			assert.Equal(t, r.URL.Path, "/report-categories")
			assert.Equal(t, r.URL.Query().Get("page"), "1")
			assert.Equal(t, r.URL.Query().Get("pageSize"), "50")
			assert.Equal(t, r.Header.Get("Authorization"), "tok_1230")
			assert.Equal(t, r.Header.Get("ST-App-Key"), "app_123")

			io.WriteString(w, `{"data": [{"name": "Category A", "id": "cat-a"}, {"name": "Category B", "id": "cat-b"}]}`)
		})

		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}
		got, err := srv.GetCategories(context.Background(), nil)
		assert.NilError(t, err)

		assert.DeepEqual(t, got, &CategoryList{
			Items: []Category{
				{
					ID:   "cat-a",
					Name: "Category A",
				},
				{
					ID:   "cat-b",
					Name: "Category B",
				},
			},
		})
	})

	t.Run("fetches specific page of categories", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			assert.NilError(t, r.ParseForm())
			assert.Equal(t, r.URL.Query().Get("page"), "4")
			assert.Equal(t, r.URL.Query().Get("pageSize"), "100")
			io.WriteString(w, `{"data": [{"name": "Category B", "id": "cat-b"}]}`)
		})

		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}
		got, err := srv.GetCategories(context.Background(), &ReportOptions{
			Page:     4,
			PageSize: 100,
		})

		assert.NilError(t, err)
		assert.DeepEqual(t, got, &CategoryList{
			Items: []Category{
				{
					ID:   "cat-b",
					Name: "Category B",
				},
			},
		})
	})

	t.Run("returns error when request fails", func(t *testing.T) {
		srv := reportService{baseURL: "", client: buildClient()}
		_, err := srv.GetCategories(context.Background(), nil)
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})

	t.Run("returns error when request building fail", func(t *testing.T) {
		srv := reportService{baseURL: string([]byte{0x7f}), client: buildClient()}
		_, err := srv.GetCategories(context.Background(), nil)
		assert.ErrorContains(t, err, "net/url: invalid control character in URL")
	})

	t.Run("returns error when non 200 response code", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "error invalid token")
		})
		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}
		_, err := srv.GetCategories(context.Background(), nil)

		want := &Error{
			StatusCode:  http.StatusUnauthorized,
			RequestPath: "/report-categories",
			Message:     "error invalid token",
		}
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns error when it fails parse json body", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "invalid json")
		})
		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}

		_, err := srv.GetCategories(context.Background(), nil)
		assert.ErrorType(t, err, &json.SyntaxError{})
	})
}

func TestReportService_GetReports(t *testing.T) {
	t.Run("fetches first page of reports for a category when report options is nil", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			assert.NilError(t, r.ParseForm())
			assert.Equal(t, r.URL.Path, "/report-category/cat-a/reports")
			assert.Equal(t, r.URL.Query().Get("page"), "1")
			assert.Equal(t, r.URL.Query().Get("pageSize"), "50")
			assert.Equal(t, r.Header.Get("Authorization"), "tok_1230")
			assert.Equal(t, r.Header.Get("ST-App-Key"), "app_123")

			io.WriteString(w, `{"data": [{"name": "Report A", "id": 1234}, {"name": "Report B", "id": 2345}]}`)
		})

		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}
		got, err := srv.GetReports(context.Background(), Category{ID: "cat-a"}, nil)
		assert.NilError(t, err)

		assert.DeepEqual(t, got, &ReportList{
			Items: []Report{
				{
					ID:   1234,
					Name: "Report A",
				},
				{
					ID:   2345,
					Name: "Report B",
				},
			},
		})
	})

	t.Run("fetches specific page of reports for a category", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			assert.NilError(t, r.ParseForm())
			assert.Equal(t, r.URL.Query().Get("page"), "5")
			assert.Equal(t, r.URL.Query().Get("pageSize"), "200")
			io.WriteString(w, `{"data": [{"name": "Report C", "id": 6666}]}`)
		})

		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}
		got, err := srv.GetReports(context.Background(), Category{ID: "cat-a"}, &ReportOptions{
			Page:     5,
			PageSize: 200,
		})

		assert.NilError(t, err)
		assert.DeepEqual(t, got, &ReportList{
			Items: []Report{
				{
					ID:   6666,
					Name: "Report C",
				},
			},
		})
	})

	t.Run("returns error when request fails", func(t *testing.T) {
		srv := reportService{baseURL: "", client: buildClient()}
		_, err := srv.GetReports(context.Background(), Category{}, nil)
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})

	t.Run("returns error when request building fail", func(t *testing.T) {
		srv := reportService{baseURL: string([]byte{0x7f}), client: buildClient()}
		_, err := srv.GetReports(context.Background(), Category{}, nil)
		assert.ErrorContains(t, err, "net/url: invalid control character in URL")
	})

	t.Run("returns error when non 200 response code", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "error invalid token")
		})
		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}
		_, err := srv.GetReports(context.Background(), Category{ID: "cat-b"}, nil)

		want := &Error{
			StatusCode:  http.StatusUnauthorized,
			RequestPath: "/report-category/cat-b/reports",
			Message:     "error invalid token",
		}
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns error when it fails parse json body", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "invalid json")
		})
		defer server.Close()

		srv := reportService{baseURL: server.URL, client: buildClient()}

		_, err := srv.GetReports(context.Background(), Category{}, nil)
		assert.ErrorType(t, err, &json.SyntaxError{})
	})
}

func buildClient() *Client {
	return &Client{
		client: http.DefaultClient,
		session: &Session{
			Token:     "tok_1230",
			ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
		},
		metadata: ClientInfo{
			AppID: "app_123",
		},
	}
}
