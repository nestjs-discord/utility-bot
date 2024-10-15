package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/nestjs-discord/utility-bot/pkg/uptime"
	"log/slog"
	"runtime"
)

func (h *Handler) MessageCreate(s *discordgo.Session, i *discordgo.MessageCreate) {
	if i.Message.Author.Bot {
		return
	}

	h.logger.Debug("message create",
		slog.String("id", i.Message.ID),
		slog.String("content", i.Message.Content),
	)

	if h.autoMod.Enabled() {
		h.autoMod.Handler(s, i)
	}

	if !h.isModerator(i.Author.ID) {
		return
	}

	switch i.Content {
	case "!stats":
		h.statsHandler(s, i)
		return
	case "!automod":
		h.autoMod.TrackHandler(s, i)
		return
	}
}

func (h *Handler) statsHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	embed := &discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Total allocated memory for heap objects",
				Value: humanize.Bytes(m.TotalAlloc),
			},
			{
				Name:  "Total memory obtained from the OS",
				Value: humanize.Bytes(m.Sys),
			},
			{
				Name:  "Allocated heap objects",
				Value: humanize.Bytes(m.Alloc),
			},
			{
				Name:  "Heap memory reserved but not allocated",
				Value: humanize.Bytes(m.HeapIdle),
			},
			{
				Name:  "Heap memory in-use",
				Value: humanize.Bytes(m.HeapInuse),
			},
			{
				Name:  "Stack memory in-use",
				Value: humanize.Bytes(m.StackInuse),
			},
			{
				Name:  "Memory obtained from system via mmap",
				Value: humanize.Bytes(m.MSpanSys + m.MCacheSys),
			},
			{
				Name:  "Memory used for GC metadata",
				Value: humanize.Bytes(m.GCSys),
			},
			{
				Name:  "Uptime",
				Value: uptime.Uptime(),
			},
		},
	}

	_, _ = s.ChannelMessageSendEmbed(i.ChannelID, embed)
}
