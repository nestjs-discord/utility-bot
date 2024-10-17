package main

import (
	"github.com/nestjs-discord/utility-bot/app"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	application, err := app.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}

	err = application.Bot.OpenWebsocketConnection()
	if err != nil {
		log.Fatalf("bot open failed: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	slog.Info("shutting down")
	err = application.Bot.Close()
	if err != nil {
		log.Fatalf("bot close failed: %v", err)
	}
	slog.Info("shutdown done")
}
