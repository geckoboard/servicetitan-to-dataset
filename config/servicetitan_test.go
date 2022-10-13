package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

func TestServiceTitan_Validate(t *testing.T) {
	t.Run("returns an error", func(t *testing.T) {
		specs := []struct {
			name    string
			in      ServiceTitan
			wantErr Error
		}{
			{
				name: "default client info returns all errors",
				in:   ServiceTitan{},
				wantErr: Error{
					scope: "servicetitan",
					messages: []string{
						"missing app_id",
						"missing tenant_id",
						"missing client_id",
						"missing client_secret",
					},
				},
			},
			{
				name: "missing app id",
				in:   ServiceTitan{TenantID: "ten_3", ClientID: "cl_14", ClientSecret: "sec_15"},
				wantErr: Error{
					scope: "servicetitan",
					messages: []string{
						"missing app_id",
					},
				},
			},
			{
				name: "missing tenant id",
				in:   ServiceTitan{AppID: "ap_3", ClientID: "cl_14", ClientSecret: "sec_15"},
				wantErr: Error{
					scope: "servicetitan",
					messages: []string{
						"missing tenant_id",
					},
				},
			},
			{
				name: "missing client id",
				in:   ServiceTitan{AppID: "ap_3", TenantID: "te_14", ClientSecret: "sec_15"},
				wantErr: Error{
					scope: "servicetitan",
					messages: []string{
						"missing client_id",
					},
				},
			},
			{
				name: "missing client secret",
				in:   ServiceTitan{AppID: "ap_3", TenantID: "te_14", ClientID: "cl_15"},
				wantErr: Error{
					scope: "servicetitan",
					messages: []string{
						"missing client_secret",
					},
				},
			},
		}

		for _, spec := range specs {
			t.Run(spec.name, func(t *testing.T) {
				assert.DeepEqual(t, spec.in.Validate(), spec.wantErr, cmp.AllowUnexported(Error{}))
			})
		}
	})

	t.Run("returns no error when all valid", func(t *testing.T) {
		in := ServiceTitan{
			AppID:        "ap_3",
			TenantID:     "te_14",
			ClientID:     "cl_15",
			ClientSecret: "sec_9",
		}

		assert.NilError(t, in.Validate())
	})
}
