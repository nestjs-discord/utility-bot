package handler

import "github.com/bwmarrin/discordgo"

func (h *Handler) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.opts.InteractionHandler.ApplicationCommand(s, i)
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.opts.InteractionHandler.ApplicationCommandAutocomplete(s, i)
		return
	case discordgo.InteractionMessageComponent:
		h.opts.InteractionHandler.MessageComponent(s, i)
		return
	case discordgo.InteractionModalSubmit:
		h.opts.InteractionHandler.ModalSubmit(s, i)
		return
	}
}
