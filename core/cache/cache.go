package cache

import "github.com/nestjs-discord/utility-bot/internal/ratelimit"

var Ratelimit *ratelimit.TTLMap

func InitRatelimit(ttl int) {
	Ratelimit = ratelimit.New(ttl)
}
