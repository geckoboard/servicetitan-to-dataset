package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestConfig_ExtractValuesFromEnv(t *testing.T) {
	os.Setenv("ENV_1", "val1")
	os.Setenv("ENV_2", "val2")
	os.Setenv("ENV_3", "val3")
	os.Setenv("ENV_4", "val4")
	os.Setenv("ENV_5", "val5")

	defer func() {
		os.Unsetenv("ENV_1")
		os.Unsetenv("ENV_2")
		os.Unsetenv("ENV_3")
		os.Unsetenv("ENV_4")
		os.Unsetenv("ENV_5")
	}()

	t.Run("replaces interpolated values from envs", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "{{ENV_1}}",
				TenantID:     "{{ENV_2}}",
				ClientSecret: "{{ENV_3}}",
				ClientID:     "{{ENV_4}}",
			},
			Geckoboard: Geckoboard{
				APIKey: "{{ENV_5}}",
			},
		}

		want := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "val1",
				TenantID:     "val2",
				ClientSecret: "val3",
				ClientID:     "val4",
			},
			Geckoboard: Geckoboard{
				APIKey: "val5",
			},
		}

		in.ExtractValuesFromEnv()
		assert.DeepEqual(t, in, want, cmpopts.IgnoreUnexported(Config{}))
	})

	t.Run("leaves un-interpolated values as is", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "{{ENV_1}}",
				TenantID:     "tenant-id",
				ClientSecret: "",
				ClientID:     "{{ENV_4}}",
			},
			Geckoboard: Geckoboard{
				APIKey: "apikey22",
			},
		}

		want := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "val1",
				TenantID:     "tenant-id",
				ClientSecret: "",
				ClientID:     "val4",
			},
			Geckoboard: Geckoboard{
				APIKey: "apikey22",
			},
		}

		in.ExtractValuesFromEnv()
		assert.DeepEqual(t, in, want, cmpopts.IgnoreUnexported(Config{}))
	})
}

func TestConfig_Validate(t *testing.T) {
	t.Run("returns errors for servicetitan", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "app",
				TenantID:     "ten",
				ClientSecret: "secret",
			},
		}

		assert.ErrorContains(t, in.Validate(), "Config section \"servicetitan\" errors:\n - missing client_id")
	})

	t.Run("returns error when invalid time location", func(t *testing.T) {
		in := Config{TimeLocation: "fake"}
		assert.ErrorContains(t, in.Validate(), "Config time_location error: unknown time zone fake")
	})

	t.Run("returns errors for geckoboard", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "app",
				TenantID:     "ten",
				ClientID:     "id",
				ClientSecret: "secret",
			},
		}

		assert.ErrorContains(t, in.Validate(), "Config section \"geckoboard\" errors:\n - missing api_key")
	})

	t.Run("returns error at least one entry required", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "app",
				TenantID:     "ten",
				ClientID:     "id",
				ClientSecret: "secret",
			},
			Geckoboard: Geckoboard{
				APIKey: "api123",
			},
		}

		assert.ErrorContains(t, in.Validate(), "Config section \"entries\" errors:\n - at least one entry is required")
	})

	t.Run("returns error when on of the entries is invalid", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "app",
				TenantID:     "ten",
				ClientID:     "id",
				ClientSecret: "secret",
			},
			Geckoboard: Geckoboard{
				APIKey: "api123",
			},
			Entries: Entries{
				{
					Dataset: Dataset{
						RequiredFields: []string{"Name"},
					},
					Report: Report{
						ID:         "rpt-1",
						CategoryID: "cat-1",
					},
				},
				{
					Dataset: Dataset{},
					Report: Report{
						ID: "rpt-1",
					},
				},
			},
		}

		assert.ErrorContains(t, in.Validate(), "Config section \"entries[2]\" errors:\n - at least one dataset required_field is required, please use the report field name as the identifier\n - category_id is required")
	})

	t.Run("allows empty time location with valid config", func(t *testing.T) {
		in := Config{
			ServiceTitan: ServiceTitan{
				AppID:        "app",
				TenantID:     "ten",
				ClientID:     "id",
				ClientSecret: "secret",
			},
			Geckoboard: Geckoboard{
				APIKey: "api123",
			},
			Entries: Entries{
				{
					Dataset: Dataset{
						RequiredFields: []string{"Name"},
					},
					Report: Report{
						ID:         "rpt-1",
						CategoryID: "cat-1",
					},
				},
			},
		}

		assert.NilError(t, in.Validate())
	})

	t.Run("returns no error with valid config", func(t *testing.T) {
		in := Config{
			TimeLocation: "America/New_York",
			ServiceTitan: ServiceTitan{
				AppID:        "app",
				TenantID:     "ten",
				ClientID:     "id",
				ClientSecret: "secret",
			},
			Geckoboard: Geckoboard{
				APIKey: "api123",
			},
			Entries: Entries{
				{
					Dataset: Dataset{
						RequiredFields: []string{"Name"},
					},
					Report: Report{
						ID:         "rpt-1",
						CategoryID: "cat-1",
					},
				},
			},
		}

		assert.NilError(t, in.Validate())
	})
}
