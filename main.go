package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	godotenv "github.com/joho/godotenv"
)

var cntAll, cntYes, cntNo int
var ages = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "60 ++"}
var ages_stat = [6][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}
var ageGroup = map[string]int{ages[0]: 15, ages[1]: 25, ages[2]: 35, ages[3]: 45, ages[4]: 55, ages[5]: 65}
var user_age = make(map[int64]int)
var user_timestamp = make(map[int64]int64)
var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"

var numericInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ages[0], ages[0]),
		tgbotapi.NewInlineKeyboardButtonData(ages[1], ages[1]),
		tgbotapi.NewInlineKeyboardButtonData(ages[2], ages[2]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ages[3], ages[3]),
		tgbotapi.NewInlineKeyboardButtonData(ages[4], ages[4]),
		tgbotapi.NewInlineKeyboardButtonData(ages[5], ages[5]),
	),
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Start"),
	),
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	API_TOKEN := os.Getenv("TELEGRAM_API_TOKEN")

	cntAll, cntYes, cntNo = 0, 0, 0

	bot, err := tgbotapi.NewBotAPI(API_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	var userID int64
	var userData string

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
			continue
		}

		userID = update.Message.From.ID
		if update.CallbackQuery != nil {
			userData = update.CallbackQuery.Data
		}
		log.Printf("userID %d, %s", userID, userData)
		if tout, ok := user_timestamp[userID]; ok {
			log.Printf("Found timestamp %d", tout)
			user_timestamp[userID] = tout
		} else {
			log.Printf("Init timestamp 0")
			user_timestamp[userID] = 10
		}

		// new user - send age question
		user_age[userID] = 0

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		log.Printf("Message text: %s", msg.Text)
		if update.Message.Text == "Start" {
			msg.ReplyMarkup = numericInlineKeyboard
		} else {
			msg.ReplyMarkup = startKeyboard
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
