package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.bot.OpenWebsocketConnection()
	if err != nil {
		log.Fatalf("bot open failed: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	slog.Info("shutting down")
	err = app.bot.Close()
	if err != nil {
		log.Fatalf("bot close failed: %v", err)
	}
	slog.Info("shutdown done")
}
