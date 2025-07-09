package ai

import (
	"context"
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"github.com/nestjs-discord/utility-bot/infra/services/gemini"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

type Forums struct {
	logger *slog.Logger
	cfg    yaml.AI
	Gemini *gemini.Gemini
}

func NewForums(cfg yaml.AI, gemini *gemini.Gemini) *Forums {
	return &Forums{
		logger: logger.NewWithSubsystem("bot", "handler", "ai", "forums"),
		cfg:    cfg,
		Gemini: gemini,
	}
}

func (f *Forums) IsChannelAllowedToUseAIFeatures(channel *dgo.Channel) (bool, error) {
	if channel.Type != dgo.ChannelTypeGuildPublicThread &&
		channel.Type != dgo.ChannelTypeGuildPrivateThread {
		return false, nil
	}

	// We only allow the first message in a thread to use AI features
	if channel.MessageCount > 0 {
		return false, nil
	}

	for _, forum := range f.cfg.Forums {
		if channel.ParentID == forum {
			return true, nil
		}
	}

	return false, nil
}

func (f *Forums) Handle(s *dgo.Session, m *dgo.MessageCreate) bool {
	if !f.cfg.Enabled || m.Author.Bot {
		return false
	}

	content := strings.TrimSpace(m.Content)
	if content == "" {
		return false
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return false
	}

	allowed, checkErr := f.IsChannelAllowedToUseAIFeatures(channel)
	if checkErr != nil || !allowed {
		return false
	}

	defer f.closeChannel(s, m)
	f.reactWithLoading(s, m)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	aiRespond, genErr := f.Gemini.Generate(ctx, channel.Name, content)
	if genErr != nil {
		f.logger.Error("error generating AI response",
			slog.String("channel_name", channel.Name),
			slog.String("channel_id", channel.ID),
			slog.Any("err", genErr),
		)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error generating AI response. Please contact a moderator.")
		f.reactWithError(s, m)
		return true
	}

	err = f.sendPromptResult(s, m.ChannelID, aiRespond)
	if err != nil {
		f.logger.Error("failed to send AI response",
			slog.String("channel_name", channel.Name),
			slog.String("channel_id", m.ChannelID),
			slog.Any("err", err),
		)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error sending AI response. Please contact a moderator.")
		f.reactWithError(s, m)
		return true
	}

	err = f.sendPromptResultSources(s, m.ChannelID, aiRespond)
	if err != nil {
		f.logger.Error("failed to send AI response sources",
			slog.String("channel_name", channel.Name),
			slog.String("channel_id", m.ChannelID),
			slog.Any("err", err),
		)
	}

	f.reactWithSuccess(s, m)

	return true
}

func (f *Forums) sendPromptResult(s *dgo.Session, channelID string, aiRespond *gemini.GenerateResponse) error {
	data := &dgo.MessageSend{}
	if len(aiRespond.Content) < 1900 {
		data.Content = aiRespond.Content
	} else {
		data.Content = "We have generated a response to your question, " +
			"but it is too long to be sent as a single message. " +
			"Please check the response attached as a file."
		data.Files = []*dgo.File{
			{
				Name:        "response.md",
				ContentType: "text/plain",
				Reader:      strings.NewReader(aiRespond.Content),
			},
		}
	}

	_, err := s.ChannelMessageSendComplex(channelID, data)
	if err != nil {
		return fmt.Errorf("ChannelMessageSendComplex: %w", err)
	}

	return nil
}

func (f *Forums) sendPromptResultSources(s *dgo.Session, channelID string, aiRespond *gemini.GenerateResponse) error {
	data := &dgo.MessageSend{}
	for _, source := range aiRespond.Sources {
		if source.Text != "" && source.URL != "" {
			data.Content += "- [" + source.Text + "](<" + source.URL + ">)\n"
		}
	}

	for _, suggestion := range aiRespond.SearchSuggestions {
		row := dgo.ActionsRow{
			Components: []dgo.MessageComponent{
				dgo.Button{
					Emoji: &dgo.ComponentEmoji{
						Name: "🔍",
					},
					Label: suggestion,
					URL:   "https://google.com/search?q=" + url.QueryEscape(suggestion),
					Style: dgo.LinkButton,
				},
			},
		}
		data.Components = append(data.Components, row)
	}

	// If there is no content and no components, we don't send a message
	if data.Content == "" && len(data.Components) == 0 {
		return nil
	}

	_, err := s.ChannelMessageSendComplex(channelID, data)
	if err != nil {
		return fmt.Errorf("ChannelMessageSendComplex: %w", err)
	}

	return nil
}

func (f *Forums) reactWithLoading(s *dgo.Session, m *dgo.MessageCreate) {
	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "⏳")
}

func (f *Forums) reactWithSuccess(s *dgo.Session, m *dgo.MessageCreate) {
	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func (f *Forums) reactWithError(s *dgo.Session, m *dgo.MessageCreate) {
	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
}

func (f *Forums) closeChannel(s *dgo.Session, msg *dgo.MessageCreate) {
	_ = s.MessageReactionsRemoveEmoji(msg.ChannelID, msg.ID, "⏳")

	archived := false
	locked := true
	_, err := s.ChannelEdit(msg.ChannelID, &dgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})
	if err != nil {
		f.logger.Error("failed to close channel",
			slog.String("channel_id", msg.ChannelID),
		)
	}
}
