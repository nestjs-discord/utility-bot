package yaml

import "errors"

type Antispam struct {
	Enabled            bool     `yaml:"enabled"`
	ModeratorsBypass   bool     `yaml:"moderatorsBypass"`
	LogChannelId       string   `yaml:"logChannelId"`
	MessageTTLSec      int      `yaml:"messageTTLSec"`
	MaxChannelsPerUser int      `yaml:"maxChannelsPerUser"`
	DenyTTLSec         int      `yaml:"denyTTLSec"`
	TrackedChannelIds  []string `yaml:"trackedChannelIds"`
}

func (a Antispam) validate() error {
	if !a.Enabled {
		return nil
	}

	if a.LogChannelId == "" {
		return errors.New("anti spam log channel id is required")
	}

	if a.MessageTTLSec == 0 {
		return errors.New("message ttl is required")
	}

	if a.MaxChannelsPerUser == 0 {
		return errors.New("max channels per user is required")
	}

	if a.DenyTTLSec == 0 {
		return errors.New("deny ttl is required")
	}

	if len(a.TrackedChannelIds) == 0 {
		return errors.New("tracked channel ids is required")
	}

	return nil
}

func NewAntispam(config *Config) Antispam {
	return config.Antispam
}
