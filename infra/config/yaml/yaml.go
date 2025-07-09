package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Path string

type Config struct {
	Moderators           Moderators           `yaml:"moderators"`
	PrivilegedForChannel PrivilegedForChannel `yaml:"privilegedForChannel"`
	Antispam             Antispam             `yaml:"antispam"`
	RateLimit            RateLimit            `yaml:"rateLimit"`
	Forms                Forms                `yaml:"forms"`
	AI                   AI                   `yaml:"ai"`
	ArchiveCommand       ArchiveCommand       `yaml:"archiveCommand"`
	SolvedCommand        SolvedCommand        `yaml:"solvedCommand"`
	Commands             Commands             `yaml:"commands"`
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
