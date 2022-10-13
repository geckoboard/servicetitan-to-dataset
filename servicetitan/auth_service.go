package servicetitan

import (
	"context"
	"net/url"
	"servicetitan-to-dataset/config"
	"strings"
	"time"
)

type authStep bool

type AuthService interface {
	GetToken(context.Context, config.ServiceTitan) (*Session, error)
}

type authService struct {
	baseURL string
	client  *Client
}

type Session struct {
	Token     string    `json:"access_token"`
	ExpireIn  int       `json:"expires_in"`
	ExpiresAt time.Time `json:"-"`
}

func (a authService) GetToken(ctx context.Context, cfg config.ServiceTitan) (*Session, error) {
	q := url.Values{}
	q.Add("grant_type", "client_credentials")
	q.Add("client_id", cfg.ClientID)
	q.Add("client_secret", cfg.ClientSecret)

	req, err := a.client.buildPOSTRequest(a.baseURL+"/connect/token", strings.NewReader(q.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	session := &Session{}

	ctx = context.WithValue(ctx, "authStep", authStep(true))
	if err := a.client.doRequest(req.WithContext(ctx), session); err != nil {
		return nil, err
	}

	session.setExpiryDate()
	return session, nil
}

func (s *Session) setExpiryDate() {
	now := time.Now().UTC()
	s.ExpiresAt = now.Add(time.Duration(s.ExpireIn) * time.Second)
}

func (s *Session) IsExpired() bool {
	// Give a minute buffer
	return time.Now().UTC().Add(time.Minute).After(s.ExpiresAt)
}
