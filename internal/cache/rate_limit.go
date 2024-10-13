package cache

import (
	"github.com/nestjs-discord/utility-bot/pkg/rate_limit"
)

var RateLimit *rate_limit.TTLMap

func InitRateLimit(ttl int) {
	RateLimit = rate_limit.New(ttl)
}
