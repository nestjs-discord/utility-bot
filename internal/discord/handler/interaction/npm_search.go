package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/command/npm"
	npmAPI "github.com/erosdesire/discord-nestjs-utility-bot/internal/npm"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/uniplaces/carbon"
	"strings"
	"time"
)

func NpmSearchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := &npmAPI.SearchOptions{
		Size: 20,
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
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if data.Total <= 0 || len(data.Objects) <= 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No results were found with the entered information.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAutoFormatHeaders(false)

	// Markdown table format
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	//table.SetCenterSeparator("|")

	//table.SetBorder(false)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	//table.SetHeader([]string{"Package name", "Version", "Last publish"})

	for _, object := range data.Objects {
		name := object.Package.Name
		version := object.Package.Version

		date := object.Package.Date

		versionColor := generateVersionColor(version)
		dateColor := generateDateColor(date)

		table.Rich([]string{name, version, humanize.Time(date)}, []tablewriter.Colors{
			{
				// tablewriter.FgWhiteColor
			},
			*versionColor,
			*dateColor,
		})
	}
	table.Render()

	content := "```ansi\n"
	content += removeTrailingWhitespace(tableString.String())
	content += "```"

	//fmt.Println(content) // TODO: remove print

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: fmt.Sprintf("View %v results on npmjs.com", humanize.Comma(data.Total)),
							Style: discordgo.LinkButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "📦",
							},
							Disabled: false,
							URL:      "https://www.npmjs.com/search?q=" + options.Text,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to respond npm-search command")

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func generateVersionColor(version string) *tablewriter.Colors {
	if strings.HasPrefix(version, "0") {
		return &tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
	}
	return &tablewriter.Colors{tablewriter.FgCyanColor}
}

func generateDateColor(date time.Time) *tablewriter.Colors {
	if carbon.Now().SubMonths(2).After(date) {
		return &tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
	}
	return &tablewriter.Colors{tablewriter.FgBlueColor}
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

func removeTrailingWhitespace(input string) string {
	lines := strings.Split(input, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.Join(lines, "\n")
}
