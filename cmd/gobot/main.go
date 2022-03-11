package main

import (
	"log"
	"os"

	"github.com/zloyboy/gobot/internal/config"
	"github.com/zloyboy/gobot/internal/database"
	"github.com/zloyboy/gobot/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	f, _ := os.OpenFile("data/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)

	cfg := config.ReadConfig()
	if cfg == nil {
		return
	}
	log.Printf("\nStart app ver %s", cfg.Version)
	dbase := database.InitDb()
	if dbase == nil {
		return
	}

	bot, err := tgbotapi.NewBotAPI(cfg.ApiToken)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true

	teleBot := telegram.NewBot(bot, dbase, cfg)
	teleBot.Run()
}
