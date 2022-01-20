package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/zloyboy/gobot/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqliteDatabase, _ := sql.Open("sqlite3", "data/stat.db")
	defer sqliteDatabase.Close()

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		return
	}
	API_TOKEN := os.Getenv("TELEGRAM_API_TOKEN")

	bot, err := tgbotapi.NewBotAPI(API_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	tBot := telegram.NewBot(bot)
	tBot.Start()
}
