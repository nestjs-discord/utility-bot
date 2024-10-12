package forms

import "github.com/bwmarrin/discordgo"

const (
	ModAcceptBtnIdPrefix = "form-mod-accept-"
	ModRejectBtnIdPrefix = "form-mod-reject-"
	ModBanBtnIdPrefix    = "form-mod-ban-"
)

func GenerateModComponents(formId string, userId string) *discordgo.ActionsRow {
	acceptBtn := discordgo.Button{
		Label:    "Accept & Publish",
		Style:    discordgo.SecondaryButton,
		Disabled: false,
		CustomID: ModAcceptBtnIdPrefix + formId,
		Emoji:    &discordgo.ComponentEmoji{Name: "✅", Animated: false},
	}
	rejectBtn := discordgo.Button{
		Label:    "Reject",
		Style:    discordgo.SecondaryButton,
		Disabled: false,
		CustomID: ModRejectBtnIdPrefix + formId,
		Emoji:    &discordgo.ComponentEmoji{Name: "❌", Animated: false},
	}
	banBtn := discordgo.Button{
		Label:    "Ban the user",
		Style:    discordgo.SecondaryButton,
		Disabled: false,
		CustomID: ModBanBtnIdPrefix + userId,
		Emoji:    &discordgo.ComponentEmoji{Name: "🔴", Animated: false},
	}

	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			acceptBtn,
			rejectBtn,
			banBtn,
		},
	}
}
