package servicetitan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client   *http.Client
	metadata ClientInfo
	session  *Session

	AuthService   AuthService
	ReportService ReportService
}

func New(info ClientInfo) (*Client, error) {
	if err := info.Validate(); err != nil {
		return nil, err
	}

	c := &Client{
		client:   &http.Client{Timeout: 30 * time.Second},
		metadata: info,
	}

	c.AuthService = authService{
		baseURL: "https://auth.servicetitan.io",
		client:  c,
	}
	c.ReportService = reportService{
		baseURL: fmt.Sprintf("https://api.servicetitan.io/reporting/v2/tenant/%s", info.TenantID),
		client:  c,
	}

	return c, nil
}

func (c *Client) buildURL(baseURL, path string, params url.Values) string {
	return fmt.Sprintf("%s?%s", baseURL+path, params.Encode())
}

func (c *Client) buildGETRequest(url string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) buildPOSTRequest(url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) addAuthorization(r *http.Request) (err error) {
	defer func() {
		r.Header.Add("Authorization", c.session.Token)
	}()

	if c.session != nil && !c.session.IsExpired() {
		return nil
	}

	c.session, err = c.AuthService.GetToken(r.Context(), c.metadata)
	return err
}

func (c *Client) doRequest(req *http.Request, resource interface{}) error {
	if authstep, _ := req.Context().Value("authStep").(authStep); !authstep {
		c.addAuthorization(req)
		req.Header.Add("ST-App-Key", c.metadata.AppID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := c.checkResponse(resp); err != nil {
		return err
	}

	if resource != nil {
		d := json.NewDecoder(resp.Body)
		if err := d.Decode(&resource); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) checkResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return &Error{
		StatusCode:  resp.StatusCode,
		RequestPath: resp.Request.URL.Path,
		Message:     string(b),
	}
}
