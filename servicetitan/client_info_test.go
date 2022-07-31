package servicetitan

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestClientInfo_Validate(t *testing.T) {
	t.Run("returns an error", func(t *testing.T) {
		specs := []struct {
			name    string
			in      ClientInfo
			wantErr error
		}{
			{
				name:    "default client info",
				in:      ClientInfo{},
				wantErr: errMissingAppID,
			},
			{
				name:    "missing app id",
				in:      ClientInfo{TenantID: "ten_3", ClientID: "cl_14", ClientSecret: "sec_15"},
				wantErr: errMissingAppID,
			},
			{
				name:    "missing tenant id",
				in:      ClientInfo{AppID: "ap_3", ClientID: "cl_14", ClientSecret: "sec_15"},
				wantErr: errMissingTenantID,
			},
			{
				name:    "missing client id",
				in:      ClientInfo{AppID: "ap_3", TenantID: "te_14", ClientSecret: "sec_15"},
				wantErr: errMissingClientID,
			},
			{
				name:    "missing client secret",
				in:      ClientInfo{AppID: "ap_3", TenantID: "te_14", ClientID: "cl_15"},
				wantErr: errMissingClientSecret,
			},
		}

		for _, spec := range specs {
			t.Run(spec.name, func(t *testing.T) {
				assert.ErrorIs(t, spec.in.Validate(), spec.wantErr)
			})
		}
	})

	t.Run("returns no error when client info valid", func(t *testing.T) {
		info := ClientInfo{
			AppID:        "ap_3",
			TenantID:     "te_14",
			ClientID:     "cl_15",
			ClientSecret: "sec_9",
		}

		assert.NilError(t, info.Validate())
	})
}
