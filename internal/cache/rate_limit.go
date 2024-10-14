package cache

import (
	"github.com/nestjs-discord/utility-bot/pkg/rate_limit"
)

var RateLimit *rate_limit.TTLMap // TODO: remove after replacing it with the one in the handler struct
