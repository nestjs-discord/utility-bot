package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func NewConfig(path Path) (*Config, error) {
	yamlFile, err := os.ReadFile(string(path))
	if err != nil {
		return nil, fmt.Errorf("unable to read the yaml config: %s", err)
	}

	var data Config
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal the yaml config: %s", err)
	}

	// TODO: manual validation

	return &data, nil
}

func NewCommands(config *Config) Commands {
	return config.Commands
}

func NewModerators(config *Config) Moderators {
	return config.Moderators
}

func NewRateLimit(config *Config) RateLimit {
	return config.RateLimit
}

func NewAntispam(config *Config) Antispam {
	return config.Antispam
}

func NewForms(config *Config) Forms {
	return config.Forms
}
