package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/nestjs-discord/utility-bot/pkg/uptime"
	"runtime"
)

func StatsHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	embed := &discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Total memory obtained from the OS",
				Value: humanize.Bytes(m.Sys),
			},
			{
				Name:  "Total allocated memory for heap objects",
				Value: humanize.Bytes(m.TotalAlloc),
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
