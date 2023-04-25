package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/command/npm"
	npmAPI "github.com/erosdesire/discord-nestjs-utility-bot/internal/npm"
	"github.com/olekukonko/tablewriter"
	"github.com/uniplaces/carbon"
	"strings"
)

func NpmSearchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := &npmAPI.SearchOptions{
		Size: 15,
	}
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case npm.SearchOptionName:
			options.Text = option.StringValue()
		case npm.SearchOptionSort:
			value := option.IntValue()
			mapSortValueToSortOptions(value, options)
		}
	}

	data, err := npmAPI.Search(options)
	if err != nil {
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	if data.Total <= 0 || len(data.Objects) <= 0 {
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "0 packages found",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	content := fmt.Sprintf("Total results: %v\n```ansi\n", data.Total)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	// Markdown table format
	table.SetAutoFormatHeaders(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	//table.SetBorder(false)
	table.SetHeader([]string{"Package name", "Version", "Last publish"})

	for _, object := range data.Objects {
		name := object.Package.Name
		version := object.Package.Version

		date := object.Package.Date

		var versionColor tablewriter.Colors
		if strings.HasPrefix(version, "0") {
			versionColor = tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
		} else {
			versionColor = tablewriter.Colors{tablewriter.FgCyanColor}
		}

		var dateColor tablewriter.Colors
		if carbon.Now().SubMonths(2).After(date) {
			dateColor = tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
		} else {
			dateColor = tablewriter.Colors{tablewriter.FgBlueColor}
		}

		table.Rich([]string{name, version, humanize.Time(date)}, []tablewriter.Colors{
			{tablewriter.FgWhiteColor},
			versionColor,
			dateColor,
		})
	}

	table.Render()

	content += tableString.String()
	content += "```"

	fmt.Println(content) // TODO: remove print

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// TODO: generate proper message
			Content: content,
		},
	})
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func mapSortValueToSortOptions(value int64, options *npmAPI.SearchOptions) {
	switch value {
	case int64(npm.SearchSortPopularity):
		options.Popularity = 1
		break
	case int64(npm.SearchSortQuality):
		options.Quality = 1
		break
	case int64(npm.SearchSortMaintenance):
		options.Maintenance = 1
		break
	}
}
