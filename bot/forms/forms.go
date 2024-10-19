package forms

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/bot/auto_mod"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"sync"
)

type Forms struct {
	cfg             yaml.Forms
	autoMod         *auto_mod.AutoMod
	markdown        *markdown.Markdown
	modActionsCache *ristretto.Cache[string, bool]
	modActionLock   sync.RWMutex
	session         *dgo.Session
}

type userInputType struct {
	InputId      string
	InputLabel   string
	InputValue   string
	AutoModCheck error
}

func NewForms(cfg yaml.Forms,
	autoMod *auto_mod.AutoMod,
	markdown *markdown.Markdown,
	session *dgo.Session,
) (*Forms, error) {
	f := &Forms{
		cfg:           cfg,
		autoMod:       autoMod,
		markdown:      markdown,
		modActionLock: sync.RWMutex{},
		session:       session,
	}

	err := f.initCacheInstance()
	if err != nil {
		return nil, err
	}

	err = f.synchronizeOpenModalButtons()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Forms) lastChannelMessage(channelId string) (*dgo.Message, error) {
	limit := 1
	beforeId := ""
	afterId := ""
	aroundId := ""

	messages, err := f.session.ChannelMessages(channelId, limit, beforeId, afterId, aroundId)
	if err != nil {
		return nil, fmt.Errorf("unable to get channel messages: %s", err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

func (f *Forms) synchronizeOpenModalButtons() error {
	for formId, form := range f.cfg {
		message, err := f.lastChannelMessage(form.ChannelId)
		if err != nil {
			return err
		}

		if message == nil {
			continue
		}

		if f.doesHaveButtonComponentWithLabel(message, form.OpenModalMessage) {
			continue
		}

		if err = f.sendOpenModalMessage(form.ChannelId, formId, form.OpenModalMessage); err != nil {
			return fmt.Errorf("unable to send form button: %s", err)
		}
	}

	return nil
}

func (f *Forms) getFormById(id string) (*yaml.Form, error) {
	form, ok := f.cfg[id]
	if !ok {
		return nil, fmt.Errorf("form '%s' not found", id)
	}

	return &form, nil
}
