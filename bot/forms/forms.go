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

type Options struct {
	Cfg      yaml.Forms
	AutoMod  *auto_mod.AutoMod
	Markdown *markdown.Markdown
	Session  *dgo.Session
}

type Forms struct {
	opts            Options
	modActionsCache *ristretto.Cache[string, bool]
	modActionLock   sync.RWMutex
}

type userInputType struct {
	InputId      string
	InputLabel   string
	InputValue   string
	AutoModCheck error
}

func NewForms(opts Options) (*Forms, error) {
	f := &Forms{
		opts:          opts,
		modActionLock: sync.RWMutex{},
	}

	err := f.initCacheInstance()
	if err != nil {
		return nil, fmt.Errorf("init form cache instance failed: %s", err)
	}

	err = f.synchronizeOpenModalButtons()
	if err != nil {
		return nil, fmt.Errorf("synchronize open modal buttons failed: %s", err)
	}

	return f, nil
}

func (f *Forms) lastChannelMessage(channelId string) (*dgo.Message, error) {
	limit := 1
	beforeId := ""
	afterId := ""
	aroundId := ""

	messages, err := f.opts.Session.ChannelMessages(channelId, limit, beforeId, afterId, aroundId)
	if err != nil {
		return nil, fmt.Errorf("unable to get channel messages: %s", err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

func (f *Forms) synchronizeOpenModalButtons() error {
	for formId, form := range f.opts.Cfg {
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
	form, ok := f.opts.Cfg[id]
	if !ok {
		return nil, fmt.Errorf("form '%s' not found", id)
	}

	return &form, nil
}
