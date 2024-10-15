package status

import "github.com/bwmarrin/discordgo"

func Update(s *discordgo.Session) error {
	activities := []*discordgo.Activity{
		{
			Name: "NestJS 😎🍿",
			Type: discordgo.ActivityTypeStreaming,
			// URL: "https://www.youtube.com/watch?v=0M8AYU_hPas",
			// URL: "https://www.twitch.tv/directory/all/tags/nestjs",

			// Revealing framework fundamentals: NestJS behind the curtain
			URL: "https://www.youtube.com/watch?v=jo-1EUxMmxc",
		},
	}

	return s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: activities,
	})
}
