package main

import (
	"log"
	"os"
	"time"

	"github.com/rewired-gh/go-ume-bot/internal/app"
	tg "gopkg.in/telebot.v3"
)

func main() {
	pref := tg.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tg.LongPoller{Timeout: 15 * time.Second},
	}

	b, err := tg.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	app.HandleCommands(b)
	b.Start()
}
