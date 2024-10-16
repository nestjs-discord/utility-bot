package forms

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"sync"
)

const (
	OpenModalButton       = "fmb"
	Modal                 = "fm"
	ModeratorAcceptButton = "fma"
	ModeratorRejectButton = "fmr"
	ModeratorBanButton    = "fmb"
)

type Forms struct {
	cfg             yaml.Forms
	modActionsCache *ristretto.Cache[string, bool]
	modActionLock   sync.RWMutex
}

func NewForms(cfg yaml.Forms, session *dgo.Session) (*Forms, error) {
	f := &Forms{
		cfg:           cfg,
		modActionLock: sync.RWMutex{},
	}

	err := f.initCacheInstance()
	if err != nil {
		return nil, err
	}

	// ----------------
	// TODO: refactor the lines below (extract method)
	// ----------------

	limit := 1
	beforeId := ""
	afterId := ""
	aroundId := ""

	for formId, form := range cfg {
		messages, err := session.ChannelMessages(form.ChannelId, limit, beforeId, afterId, aroundId)
		if err != nil {
			return nil, err
		}

		if len(messages) != 0 && f.doesHaveButtonComponentWithLabel(messages[0], form.ButtonLabel) {
			continue
		}

		label := form.ButtonLabel
		if err = f.sendFormButton(session, form.ChannelId, formId, label); err != nil {
			return nil, err
		}
	}

	return f, nil
}

type UserInput struct { // TODO: make this private?
	InputId string
	Value   string
}
