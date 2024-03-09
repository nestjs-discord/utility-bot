package automod

import (
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/config"
	"sync"
	"time"
)

type (
	UserId    string
	ChannelId string
)

type AutoMod struct {
	sync               sync.RWMutex
	userMap            map[UserId]map[ChannelId]Message
	trackedChannelsIds []ChannelId
	denyTTL            time.Duration
	deniedList         *ristretto.Cache
}

type Option struct {
	MessageTTL int
	DenyTTL    int
}

func NewAutoMod(opt Option) *AutoMod {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	a := &AutoMod{
		sync:       sync.RWMutex{},
		userMap:    make(map[UserId]map[ChannelId]Message, 0),
		denyTTL:    time.Duration(opt.DenyTTL) * time.Second,
		deniedList: cache,
	}

	go a.backgroundCleaner(opt.MessageTTL)

	return a
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
	return a.getChannelsLengthByUserId(userId) <= config.GetYaml().AutoMod.MaxChannelsLimitPerUser
}
