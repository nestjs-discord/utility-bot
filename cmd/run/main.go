package main

import (
	"github.com/nestjs-discord/utility-bot/infra/ioc"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	application, cleanup, err := ioc.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	err = application.Bot.Open()
	if err != nil {
		log.Fatalf("bot open failed: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop
	slog.Info("shutting down")
}
