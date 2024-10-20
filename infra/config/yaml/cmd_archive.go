package yaml

import "errors"

type ArchiveCommand struct {
	Description string `yaml:"description"`
	Response    string `yaml:"response"`
}

func (a ArchiveCommand) validate() error {
	if a.Description == "" {
		return errors.New("description is required")
	}
	if len(a.Description) > 60 {
		return errors.New("description is too long")
	}

	if a.Response == "" {
		return errors.New("response is required")
	}
	if len(a.Response) > 1000 {
		return errors.New("response is too long")
	}

	return nil
}

func NewArchiveCommand(config *Config) ArchiveCommand {
	return config.ArchiveCommand
}
