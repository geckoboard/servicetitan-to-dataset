package config

type Geckoboard struct {
	APIKey string `yaml:"api_key"`
}

func (gb *Geckoboard) Validate() error {
	if gb.APIKey == "" {
		return Error{
			scope:    "geckoboard",
			messages: []string{"missing api_key"},
		}
	}

	return nil
}

func (gb *Geckoboard) replaceInterpolatedValues() {
	gb.APIKey = convertEnvToValue(gb.APIKey)
}
