package automod

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/bot/moderator"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
	"sync"
	"time"
)

type (
	UserId string // TODO: remove this type
)

type AutoMod struct { // TODO: rename to antispam
	logger     *slog.Logger
	cfg        yaml.AutoMod
	sync       sync.RWMutex
	userMap    map[UserId]map[string]Message
	denyTTL    time.Duration
	deniedList *ristretto.Cache[string, bool]
	moderator  *moderator.Moderator
}

func NewAutoMod(cfg yaml.AutoMod, moderator *moderator.Moderator) (*AutoMod, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, bool]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init automod cache: %s", err)
	}

	a := &AutoMod{
		logger:     logger.NewWithSubsystem("bot", "antispam"),
		cfg:        cfg,
		sync:       sync.RWMutex{},
		userMap:    make(map[UserId]map[string]Message),
		denyTTL:    time.Duration(cfg.DenyTTL) * time.Second,
		deniedList: cache,
		moderator:  moderator,
	}

	go a.backgroundCleaner(cfg.MessageTTL)

	return a, nil
}

func (a *AutoMod) Enabled() bool {
	return a.cfg.Enabled
}

func (a *AutoMod) backgroundCleaner(ttl int) {
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

func (a *AutoMod) getChannelsLengthByUserId(id UserId) int {
	a.sync.Lock()
	defer a.sync.Unlock()

	if v, ok := a.userMap[id]; ok {
		return len(v)
	}

	return 0
}

func (a *AutoMod) IsUserWithinMaxChannelsLimit(userId UserId) bool {
	return a.getChannelsLengthByUserId(userId) <= a.cfg.MaxChannelsLimitPerUser
}
