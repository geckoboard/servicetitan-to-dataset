package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	yaml "gopkg.in/yaml.v3"
)

var interpolateRegex = regexp.MustCompile(`{{\s*([a-zA-Z0-9_]+)\s*}}`)

type Config struct {
	TimeLocation   string       `yaml:"time_location"`
	ServiceTitan   ServiceTitan `yaml:"servicetitan"`
	Geckoboard     Geckoboard   `yaml:"geckoboard"`
	RefreshTimeSec int          `yaml:"refresh_time"`
	Entries        Entries      `yaml:"entries"`

	cachedTimeLocation *time.Location
}

func LoadFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	conf := &Config{}
	if err := yaml.NewDecoder(f).Decode(conf); err != nil {
		return nil, fmt.Errorf("%s: %w", "Reading file contents failed", err)
	}

	conf.ExtractValuesFromEnv()
	return conf, nil
}

func (c *Config) ExtractValuesFromEnv() {
	c.ServiceTitan.replaceInterpolatedValues()
	c.Geckoboard.replaceInterpolatedValues()
}

func (c *Config) Validate() error {
	if err := c.loadAndValidateTimeLocation(); err != nil {
		return fmt.Errorf("Config time_location error: %v", err)
	}

	if err := c.ServiceTitan.Validate(); err != nil {
		return err
	}

	if err := c.Geckoboard.Validate(); err != nil {
		return err
	}

	if err := c.Entries.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *Config) TimeLoc() *time.Location {
	return c.cachedTimeLocation
}

func (c *Config) loadAndValidateTimeLocation() error {
	if c.TimeLocation == "" {
		return nil
	}

	loc, err := time.LoadLocation(c.TimeLocation)
	if err != nil {
		return err
	}

	c.cachedTimeLocation = loc
	return nil
}

func convertEnvToValue(value string) string {
	if value == "" {
		return ""
	}

	keys := interpolateRegex.FindStringSubmatch(value)

	if len(keys) != 2 {
		return value
	}

	return os.Getenv(keys[1])
}
