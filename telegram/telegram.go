package telegram

import (
	"log"
	"time"

	"github.com/zloyboy/gobot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ages = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "60 ++"}
var ages_stat = [6][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}
var ageGroup = map[string]int{ages[0]: 15, ages[1]: 25, ages[2]: 35, ages[3]: 45, ages[4]: 55, ages[5]: 65}
var user_age = make(map[int64]int)
var user_timestamp = make(map[int64]int64)
var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
var start_msg = "Независимый подсчет статистики по COVID-19\nУкажите вашу возрастную группу:"

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

type Bot struct {
	bot   *tgbotapi.BotAPI
	dbase *database.Dbase
}

func NewBot(bot *tgbotapi.BotAPI, db *database.Dbase) *Bot {
	return &Bot{bot: bot, dbase: db}
}

func (b *Bot) stillTimeout(userID int64) bool {
	curr_time := time.Now().Unix()
	if tstamp, ok := user_timestamp[userID]; ok {
		log.Printf("Timestamp %d, pass %d", tstamp, (curr_time - tstamp))
		if (curr_time - tstamp) < 10 {
			return true
		}
	}
	user_timestamp[userID] = curr_time
	return false
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	var chatID int64
	var userID int64
	var userData string

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
			continue
		}

		userID = update.SentFrom().ID
		chatID = update.FromChat().ID
		if update.Message != nil {
			userData = update.Message.Text
		}
		if update.CallbackQuery != nil {
			userData = update.CallbackQuery.Data
		}
		log.Printf("user %d, data %s", userID, userData)

		msg := tgbotapi.NewMessage(chatID, "")

		// first call from user
		if update.Message != nil {
			// repeat start timeout is 10 sec
			if b.stillTimeout(userID) {
				continue
			}

			if b.dbase.CheckIdName(userID) {
				// exist user - send statistic
			} else {
				// new user - send age question
				user_age[userID] = 0
				msg.ReplyMarkup = numericInlineKeyboard
				msg.Text = start_msg
			}
		} else {
			msg.ReplyMarkup = startKeyboard
			msg.Text = repeat_msg
		}

		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

	return nil
}
