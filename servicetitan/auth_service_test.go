package servicetitan

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"servicetitan-to-dataset/config"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

var (
	validClientOptions = func() config.ServiceTitan {
		return config.ServiceTitan{
			AppID:        "app_123",
			TenantID:     "tenant_123",
			ClientID:     "client_123",
			ClientSecret: "secret_123",
		}
	}
)

func TestAuthService_GetToken(t *testing.T) {
	t.Run("returns a valid session", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			assert.NilError(t, r.ParseForm())
			assert.Equal(t, r.URL.Path, "/connect/token")
			assert.Equal(t, r.Form.Get("grant_type"), "client_credentials")
			assert.Equal(t, r.Form.Get("client_id"), "client_123")
			assert.Equal(t, r.Form.Get("client_secret"), "secret_123")

			io.WriteString(w, `{"access_token":"tok_123","expires_in":600}`)
		})

		defer server.Close()

		auth := authService{baseURL: server.URL, client: &Client{client: http.DefaultClient}}
		session, err := auth.GetToken(context.Background(), validClientOptions())
		assert.NilError(t, err)

		assert.Equal(t, session.Token, "tok_123")
		assert.Equal(t, session.ExpireIn, 600)
		assert.DeepEqual(t,
			session.ExpiresAt,
			time.Now().UTC().Add(10*time.Minute),
			cmpopts.EquateApproxTime(time.Second),
		)
	})

	t.Run("returns error when request fails", func(t *testing.T) {
		auth := authService{baseURL: "", client: &Client{client: http.DefaultClient}}
		_, err := auth.GetToken(context.Background(), validClientOptions())
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})

	t.Run("returns error when request building fail", func(t *testing.T) {
		auth := authService{baseURL: string([]byte{0x7f}), client: &Client{}}
		_, err := auth.GetToken(context.Background(), validClientOptions())
		assert.ErrorContains(t, err, "net/url: invalid control character in URL")
	})

	t.Run("returns error when non 200 response code", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "error invalid client id")
		})
		defer server.Close()

		auth := authService{baseURL: server.URL, client: &Client{client: http.DefaultClient}}
		_, err := auth.GetToken(context.Background(), validClientOptions())

		want := &Error{
			StatusCode:  http.StatusForbidden,
			RequestPath: "/connect/token",
			Message:     "error invalid client id",
		}
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns error when it fails parse json body", func(t *testing.T) {
		server := buildMockServer(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "invalid json")
		})
		defer server.Close()

		auth := authService{baseURL: server.URL, client: &Client{client: http.DefaultClient}}

		_, err := auth.GetToken(context.Background(), validClientOptions())
		assert.ErrorType(t, err, &json.SyntaxError{})
	})
}

func TestSession_IsExpired(t *testing.T) {
	t.Run("returns true when now is after session expiry", func(t *testing.T) {
		session := Session{ExpiresAt: time.Now().UTC().Add(-time.Minute)}
		assert.Assert(t, session.IsExpired())
	})

	t.Run("returns true when expires is within the next minute", func(t *testing.T) {
		session := Session{ExpiresAt: time.Now().UTC().Add(59 * time.Second)}
		assert.Assert(t, session.IsExpired())
	})

	t.Run("returns false when now is before session expiry", func(t *testing.T) {
		session := Session{ExpiresAt: time.Now().UTC().Add(5 * time.Minute)}
		assert.Assert(t, !session.IsExpired())
	})
}

func buildMockServer(handlerFn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handlerFn))
}
