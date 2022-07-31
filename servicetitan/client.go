package servicetitan

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client   *http.Client
	metadata ClientInfo
	session  *Session

	AuthService AuthService
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

	return c, nil
}

func (c *Client) buildPOSTRequest(url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) doRequest(req *http.Request, resource interface{}) error {
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
