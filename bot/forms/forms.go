package forms

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/config/yaml"
)

const (
	FormButtonIdPrefix   = "form-"
	FormModalIdPrefix    = "modal-"
	ModAcceptBtnIdPrefix = "form-mod-accept-"
	ModRejectBtnIdPrefix = "form-mod-reject-"
	ModBanBtnIdPrefix    = "form-mod-ban-"
)

type Forms struct {
	cfg             yaml.Forms
	modActionsCache *ristretto.Cache[string, bool]
}

func NewForms(cfg yaml.Forms, session *discordgo.Session) (*Forms, error) {
	forms := &Forms{
		cfg: cfg,
	}

	cache, err := ristretto.NewCache(&ristretto.Config[string, bool]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cache instance for mod actions: %w", err)
	}
	forms.modActionsCache = cache

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

		if len(messages) != 0 && forms.doesHaveButtonComponentWithLabel(messages[0], form.ButtonLabel) {
			continue
		}

		label := form.ButtonLabel
		if err = forms.sendFormButton(session, form.ChannelId, formId, label); err != nil {
			return nil, err
		}
	}

	return forms, nil
}

type UserInput struct { // TODO: make this private?
	InputId string
	Value   string
}
