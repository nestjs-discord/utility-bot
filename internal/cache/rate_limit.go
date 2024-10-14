package cache

import (
	"github.com/nestjs-discord/utility-bot/config/yaml"
	"github.com/nestjs-discord/utility-bot/pkg/rate_limit"
)

var RateLimit *rate_limit.TTLMap

func Initialize(cfg yaml.RateLimit) {
	RateLimit = rate_limit.New(cfg.TTL)
}
