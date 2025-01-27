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
	opts              Options
	logger            *slog.Logger
	sync              sync.RWMutex
	userToMessagesMap map[userIdType][]Message
	denyTTL           time.Duration
	deniedList        *ristretto.Cache[string, bool]
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
		logger:            logger.NewWithSubsystem("bot", "antispam"),
		opts:              opts,
		sync:              sync.RWMutex{},
		userToMessagesMap: make(map[userIdType][]Message),
		denyTTL:           time.Duration(opts.Cfg.DenyTTLSec) * time.Second,
		deniedList:        cache,
	}

	go a.backgroundCleaner(opts.Cfg.MessageTTLSec)

	return a, nil
}

func (a *Antispam) Enabled() bool {
	return a.opts.Cfg.Enabled
}

// Removes slice element at index(s) and returns new slice
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func (a *Antispam) backgroundCleaner(ttl int) {
	for now := range time.Tick(time.Second) {
		a.sync.Lock()

		for uId := range a.userToMessagesMap {
			// Remove user id from the map if it doesn't have any messages
			if len(a.userToMessagesMap[uId]) == 0 {
				delete(a.userToMessagesMap, uId)
				break
			}

			// Remove expired messages
			for msgIndex, msg := range a.userToMessagesMap[uId] {
				if now.UTC().Unix()-msg.CreatedAt > int64(ttl) {
					a.userToMessagesMap[uId] = remove(a.userToMessagesMap[uId], msgIndex)
				}
			}
		}

		a.sync.Unlock()
	}
}

func (a *Antispam) getChannelsLengthByUserId(id userIdType) int {
	a.sync.Lock()
	defer a.sync.Unlock()

	uniqueChannelsMap := map[string]bool{}

	for _, msg := range a.userToMessagesMap[id] {
		uniqueChannelsMap[msg.ChannelID] = true
	}

	return len(uniqueChannelsMap)
}

func (a *Antispam) IsUserWithinMaxChannelsLimit(userId userIdType) bool {
	return a.getChannelsLengthByUserId(userId) <= a.opts.Cfg.MaxChannelsPerUser
}

func (a *Antispam) GetUserRepeatedMessages(userId userIdType) []Message {
	a.sync.Lock()
	defer a.sync.Unlock()

	// Result slice to store the earliest duplicate messages
	var duplicates []Message

	// Check if the user has any messages in the map
	userMessages, exists := a.userToMessagesMap[userId]
	if !exists {
		return duplicates
	}

	// Map to track message content and their occurrences
	messageContentMap := make(map[string][]Message)

	// Iterate over user's messages and group them by content
	for _, msg := range userMessages {
		trackBy := msg.ChannelID + msg.Content
		messageContentMap[trackBy] = append(messageContentMap[trackBy], msg)
	}

	// Identify duplicates
	maxAllowedRepeatCount := 3
	for _, messages := range messageContentMap {
		if len(messages) < maxAllowedRepeatCount {
			continue
		}

		duplicates = append(duplicates, messages[0])
	}

	return duplicates
}
