package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/rewired-gh/go-ume-bot/internal/app"
	"github.com/rewired-gh/go-ume-bot/internal/util"
	tg "gopkg.in/telebot.v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	config, err := util.LoadConfig("./")
	if err != nil {
		panic(err)
	}
	pref := tg.Settings{
		Token:  config.Token,
		Poller: &tg.LongPoller{Timeout: 15 * time.Second},
	}
	b, err := tg.NewBot(pref)
	if err != nil {
		slog.Error("Failed to create Telegram bot", "error", err)
		os.Exit(1)
		return
	}
	app.HandleCommands(b, config)
	b.Start()
}
