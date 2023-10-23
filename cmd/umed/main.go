package main

import (
	"log"
	"time"

	"github.com/rewired-gh/go-ume-bot/internal/app"
	"github.com/rewired-gh/go-ume-bot/internal/util"
	tg "gopkg.in/telebot.v3"
)

func main() {
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
		log.Fatal(err)
		return
	}

	app.HandleCommands(b, config)
	b.Start()
}
