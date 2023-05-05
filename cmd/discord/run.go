package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/core/cache"
	"github.com/nestjs-discord/utility-bot/core/config"
	internalDiscord "github.com/nestjs-discord/utility-bot/internal/discord"
	"github.com/nestjs-discord/utility-bot/internal/discord/command"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler"
	"github.com/nestjs-discord/utility-bot/internal/prometheus"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	metricsAddr string
)

func init() {
	Run.PersistentFlags().StringVarP(&metricsAddr,
		"metrics-addr", "", "0.0.0.0:2112", "address to run the HTTP metrics server")
}

var Run = &cobra.Command{
	Use:   "discord:run",
	Short: "Starts the Discord bot",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if metricsAddr == "" {
			return errors.New("metrics-addr flag is required.")
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Cache Markdown content
		err := cache.Content()
		if err != nil {
			return err
		}

		cache.InitRatelimit(config.GetConfig().Ratelimit.TTL)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := internalDiscord.NewSession()
		if err != nil {
			return err
		}

		c, err := session.ApplicationCommands(config.GetAppID(), config.GetGuildID())
		if err != nil {
			return errors.Wrap(err, "failed to fetch registered application commands")
		}
		command.RegisteredCommands = append(command.RegisteredCommands, c...)

		command.GenerateRegisteredCommandsSlice()

		command.RegisterStaticCommands(session)
		command.RegisterContentCommands(session)

		// Discord event handlers
		session.AddHandler(handler.Ready)

		//session.AddHandler(handlers.MessageCreate)
		session.AddHandler(handler.InteractionCreate)

		// We only care about receiving message events
		session.Identify.Intents = config.BotIntents

		// Open a websocket connection to Discord and begin listening
		err = session.Open()
		if err != nil {
			return fmt.Errorf("failed to open Discord connection: %v", err)
		}

		initPrometheus(session)

		// Graceful shutdown
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
		signal := <-sc

		log.Warn().Str("signal", signal.String()).Msg("shutting down")

		// Cleanly close down the Discord session
		return session.Close()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Warn().Msg("discord session closed")
	},
}

func initPrometheus(s *discordgo.Session) {
	p := prometheus.New()

	go func() {
		for range time.Tick(10 * time.Second) {
			value := s.HeartbeatLatency().Seconds()
			p.SetHeartbeatLatency(value)
		}
	}()

	go func() {
		err := p.ListenAndServe(metricsAddr)
		if err != nil {
			log.Panic().Err(err).Msg("failed to listen http server")
		}
	}()
}
