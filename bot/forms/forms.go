package forms

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/ristretto"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"sync"
)

type Forms struct {
	cfg             yaml.Forms
	modActionsCache *ristretto.Cache[string, bool]
	modActionLock   sync.RWMutex
}

type userInput struct {
	InputId string
	Value   string
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

	err = f.synchronizeOpenModalButtons(session)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Forms) synchronizeOpenModalButtons(session *dgo.Session) error {
	for formId, form := range f.cfg {
		message, err := f.getLastChannelMessage(session, form.ChannelId)
		if err != nil {
			return err
		}

		if message == nil {
			continue
		}

		if f.doesHaveButtonComponentWithLabel(message, form.ButtonLabel) {
			continue
		}

		label := form.ButtonLabel
		if err = f.sendFormButton(session, form.ChannelId, formId, label); err != nil {
			return fmt.Errorf("unable to send form button: %s", err)
		}
	}

	return nil
}

func (f *Forms) getLastChannelMessage(session *dgo.Session, channelId string) (*dgo.Message, error) {
	limit := 1
	beforeId := ""
	afterId := ""
	aroundId := ""

	messages, err := session.ChannelMessages(channelId, limit, beforeId, afterId, aroundId)
	if err != nil {
		return nil, fmt.Errorf("unable to get channel messages: %s", err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

func (f *Forms) getFormById(id string) (*yaml.Form, error) {
	form, ok := f.cfg[id]
	if !ok {
		return nil, fmt.Errorf("form '%s' not found", id)
	}

	return &form, nil
}
