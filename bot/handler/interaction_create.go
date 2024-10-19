package handler

import (
	dgo "github.com/bwmarrin/discordgo"
)

func (h *Handler) InteractionCreate(s *dgo.Session, i *dgo.InteractionCreate) {
	switch i.Type {
	case dgo.InteractionApplicationCommand:
		h.opts.InteractionHandler.ApplicationCommand(s, i)
		return
	case dgo.InteractionApplicationCommandAutocomplete:
		h.opts.InteractionHandler.ApplicationCommandAutocomplete(s, i)
		return
	case dgo.InteractionMessageComponent:
		h.opts.InteractionHandler.MessageComponent(s, i)
		return
	case dgo.InteractionModalSubmit:
		h.opts.InteractionHandler.ModalSubmit(s, i)
		return
	}
}
