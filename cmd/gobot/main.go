package main

import (
	"log"
	"os"

	"github.com/zloyboy/gobot/internal/database"
	"github.com/zloyboy/gobot/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

func main() {
	f, _ := os.OpenFile("data/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)

	dbase := database.InitDb()
	if dbase == nil {
		return
	}
	if err := godotenv.Load("data/.env"); err != nil {
		log.Print("No .env file found")
		return
	}
	API_TOKEN := os.Getenv("TELEGRAM_API_TOKEN")

	bot, err := tgbotapi.NewBotAPI(API_TOKEN)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true

	teleBot := telegram.NewBot(bot, dbase)
	teleBot.Run()
}
