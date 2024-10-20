package yaml

import "errors"

type RateLimit struct {
	TTLSec   int    `yaml:"ttlSec"`
	MaxUsage int    `yaml:"maxUsage"`
	Message  string `yaml:"message"`
}

func (r RateLimit) validate() error {
	if r.TTLSec == 0 {
		return errors.New("ttlSec must be greater than zero")
	}
	if r.MaxUsage < 1 {
		return errors.New("maxUsage must be greater than one")
	}
	if r.Message == "" {
		return errors.New("message is required")
	}
	if len(r.Message) < 3 {
		return errors.New("message is too short")
	}

	return nil
}

func NewRateLimit(config *Config) RateLimit {
	return config.RateLimit
}
