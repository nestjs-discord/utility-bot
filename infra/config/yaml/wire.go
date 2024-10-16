package yaml

import "github.com/google/wire"

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

var Set = wire.NewSet(
	wire.NewSet(
		wire.Value(Path("config.yml")),
		NewConfig,
	),
	wire.NewSet(
		NewCommands,
		NewModerators,
		NewRateLimit,
		NewAntispam,
		NewForms,
	),
)
