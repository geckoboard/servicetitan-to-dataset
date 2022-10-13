package config

type ServiceTitan struct {
	AppID        string `yaml:"app_id"`
	TenantID     string `yaml:"tenant_id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

func (st *ServiceTitan) Validate() error {
	var msgs []string

	if st.AppID == "" {
		msgs = append(msgs, "missing app_id")
	}

	if st.TenantID == "" {
		msgs = append(msgs, "missing tenant_id")
	}

	if st.ClientID == "" {
		msgs = append(msgs, "missing client_id")
	}

	if st.ClientSecret == "" {
		msgs = append(msgs, "missing client_secret")
	}

	if len(msgs) > 0 {
		return Error{
			scope:    "servicetitan",
			messages: msgs,
		}
	}

	return nil
}

func (st *ServiceTitan) replaceInterpolatedValues() {
	st.AppID = convertEnvToValue(st.AppID)
	st.TenantID = convertEnvToValue(st.TenantID)
	st.ClientID = convertEnvToValue(st.ClientID)
	st.ClientSecret = convertEnvToValue(st.ClientSecret)
}
