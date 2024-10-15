package rate_limit

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/moderator"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"sync"
	"time"
)

// RateLimit is a thread-safe map-based implementation of a TTL (Time to Live) rate limiter
// It allows you to track usage count for each key and automatically evict stale entries from the map.
//
// Usage:
//
//	r := NewRateLimit(60) // Creates a new RateLimit with max TTL of 60 seconds
//	r.IncrementUsage("key") // Increments usage count for "key"
//	count := r.GetUsageCount("key") // Gets the current usage count for "key"
//
// The RateLimit instance created with New() will automatically evict entries that have not been accessed for more than
// the specified TTL. This eviction process is done asynchronously by a goroutine.
type RateLimit struct {
	cfg       yaml.RateLimit
	moderator *moderator.Moderator
	m         map[string]*item // The underlying map that holds the key-value pairs
	l         sync.Mutex       // The mutex used to synchronize access to the map
}

// item is a struct that represents a value in the map along with its creation timestamp
type item struct {
	value     int   // The usage count for this item
	createdTs int64 // The UNIX timestamp when this item was created
}

// NewRateLimit returns a new RateLimit instance with a maximum TTL of maxTTL seconds.
// The returned instance automatically evicts stale entries every second.
func NewRateLimit(cfg yaml.RateLimit, moderator *moderator.Moderator) *RateLimit {
	r := &RateLimit{
		cfg:       cfg,
		moderator: moderator,
		m:         make(map[string]*item),
	}

	go func() {
		for now := range time.Tick(time.Second) {
			r.l.Lock()
			for k, v := range r.m {
				if now.Unix()-v.createdTs > int64(cfg.TTL) {
					delete(r.m, k)
				}
			}
			r.l.Unlock()
		}
	}()

	return r
}

// IncrementUsage increments the usage count for the specified key.
// If the key does not exist in the map, it is added with a usage count of 1.
func (r *RateLimit) IncrementUsage(k string) {
	r.l.Lock()
	defer r.l.Unlock()

	_, ok := r.m[k]
	if !ok {
		r.m[k] = &item{
			value:     1,
			createdTs: time.Now().Unix(),
		}
	} else {
		r.m[k].value++
	}
}

// GetUsageCount returns the current usage counts for the specified key.
// If the key does not exist in the map, it returns 0.
func (r *RateLimit) GetUsageCount(k string) (v int) {
	r.l.Lock()
	defer r.l.Unlock()

	if it, ok := r.m[k]; ok {
		v = it.value
	}

	return v
}

func (r *RateLimit) CheckRateLimit(userID string) bool {
	if r.moderator.IsUserModerator(userID) {
		return false
	}

	r.IncrementUsage(userID)

	return r.GetUsageCount(userID) > r.cfg.Usage
}

func (r *RateLimit) ForbidInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: r.cfg.Message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
