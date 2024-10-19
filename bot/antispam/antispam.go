package antispam

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
	"sync"
	"time"
)

type (
	userIdType string
)

type Options struct {
	Cfg        yaml.Antispam
	Moderators *moderators.Moderators
}

type Antispam struct {
	opts       Options
	logger     *slog.Logger
	sync       sync.RWMutex
	userMap    map[userIdType]map[string]Message
	denyTTL    time.Duration
	deniedList *ristretto.Cache[string, bool]
}

func NewAntispam(opts Options) (*Antispam, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, bool]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init antispam cache: %s", err)
	}

	a := &Antispam{
		logger:     logger.NewWithSubsystem("bot", "antispam"),
		opts:       opts,
		sync:       sync.RWMutex{},
		userMap:    make(map[userIdType]map[string]Message),
		denyTTL:    time.Duration(opts.Cfg.DenyTTLSec) * time.Second,
		deniedList: cache,
	}

	go a.backgroundCleaner(opts.Cfg.MessageTTLSec)

	return a, nil
}

func (a *Antispam) Enabled() bool {
	return a.opts.Cfg.Enabled
}

func (a *Antispam) backgroundCleaner(ttl int) {
	for now := range time.Tick(time.Second) {
		a.sync.Lock()

		for uId := range a.userMap {
			// Remove user id from the map if it doesn't have any channel
			if len(a.userMap[uId]) == 0 {
				delete(a.userMap, uId)
				continue
			}

			// Remove expired channel ids
			for cId := range a.userMap[uId] {
				if now.UTC().Unix()-a.userMap[uId][cId].CreatedAt > int64(ttl) {
					delete(a.userMap[uId], cId)
				}
			}
		}

		a.sync.Unlock()
	}
}

func (a *Antispam) getChannelsLengthByUserId(id userIdType) int {
	a.sync.Lock()
	defer a.sync.Unlock()

	if v, ok := a.userMap[id]; ok {
		return len(v)
	}

	return 0
}

func (a *Antispam) IsUserWithinMaxChannelsLimit(userId userIdType) bool {
	return a.getChannelsLengthByUserId(userId) <= a.opts.Cfg.MaxChannelsPerUser
}
