package stats

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/nestjs-discord/utility-bot/internal/uptime"
	"runtime"
)

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
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
				},
			},
		},
	})
	if err != nil {
		util.InteractionRespondError(err, s, i)
	}
}
