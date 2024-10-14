package automod

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/config/yaml"
	"sync"
	"time"
)

type (
	UserId    string
	ChannelId string
)

type AutoMod struct {
	cfg                yaml.AutoMod
	sync               sync.RWMutex
	userMap            map[UserId]map[ChannelId]Message
	trackedChannelsIds []ChannelId // TODO: remove this in favour of the `cfg.channelIds` array
	denyTTL            time.Duration
	deniedList         *ristretto.Cache[string, bool]
}

func NewAutoMod(cfg yaml.AutoMod) (*AutoMod, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, bool]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init automod cache: %s", err)
	}

	a := &AutoMod{
		cfg:        cfg,
		sync:       sync.RWMutex{},
		userMap:    make(map[UserId]map[ChannelId]Message, 0),
		denyTTL:    time.Duration(cfg.DenyTTL) * time.Second,
		deniedList: cache,
	}

	go a.backgroundCleaner(cfg.MessageTTL)

	a.setChannels(cfg.ChannelIds)

	return a, nil
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
