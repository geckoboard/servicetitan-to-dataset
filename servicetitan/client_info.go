package servicetitan

import "errors"

var (
	errMissingAppID        = errors.New("missing app id")
	errMissingTenantID     = errors.New("missing tenant id")
	errMissingClientID     = errors.New("missing client id")
	errMissingClientSecret = errors.New("missing client secret")
)

type ClientInfo struct {
	AppID        string
	TenantID     string
	ClientID     string
	ClientSecret string
}

func (c ClientInfo) Validate() error {
	if c.AppID == "" {
		return errMissingAppID
	}

	if c.TenantID == "" {
		return errMissingTenantID
	}

	if c.ClientID == "" {
		return errMissingClientID
	}

	if c.ClientSecret == "" {
		return errMissingClientSecret
	}

	return nil
}
