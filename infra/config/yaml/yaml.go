package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// TODO: remove all the "validate" tags and perform manual validation

type Path string

type Config struct {
	Moderators Moderators `yaml:"moderators"`
	Antispam   Antispam   `yaml:"antispam"`
	RateLimit  RateLimit  `yaml:"rateLimit"`
	Forms      Forms      `yaml:"forms"`
	Commands   Commands   `yaml:"commands"`
}

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

	err = data.validate()
	if err != nil {
		return nil, fmt.Errorf("yaml validation failed: %s", err)
	}

	return &data, nil
}
