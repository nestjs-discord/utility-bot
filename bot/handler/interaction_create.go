package handler

import "github.com/bwmarrin/discordgo"

func (h *Handler) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.interactionHandler.ApplicationCommand(s, i)
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.interactionHandler.ApplicationCommandAutocomplete(s, i)
		return
	case discordgo.InteractionMessageComponent:
		h.interactionHandler.MessageComponent(s, i)
		return
	case discordgo.InteractionModalSubmit:
		h.interactionHandler.ModalSubmit(s, i)
		return
	}
}
