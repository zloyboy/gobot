package telegram

import (
	"log"

	"github.com/zloyboy/gobot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot   *tgbotapi.BotAPI
	dbase *database.Dbase
	stat  *Static
	utime *TimeStamp
}

func NewBot(bot *tgbotapi.BotAPI, db *database.Dbase) *Bot {
	return &Bot{bot: bot, dbase: db, stat: Stat(), utime: MakeStamp()}
}

func (b *Bot) Run() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	if !b.readCountryFromDb() {
		return
	}
	b.readStatFromDb()

	go b.utime.DeleteTimeouts()

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
			continue
		}

		userID := update.SentFrom().ID
		if b.utime.UserTimeout(userID) {
			continue
		}

		var userData string
		if update.Message != nil {
			userData = update.Message.Text
		} else if update.CallbackQuery != nil {
			userData = update.CallbackQuery.Data
		}

		if _, ok := user_session[userID]; ok {
			log.Printf("doing user %d, data %s", userID, userData)
			user_session[userID].userChan <- userData
		} else {
			chatID := update.FromChat().ID
			userName := update.SentFrom().FirstName
			log.Printf("Start user %d, data %s", userID, userData)

			user_session[userID] = MakeSession(b, userID, chatID, userName)
			go user_session[userID].RunSurvey(user_session[userID].userChan)
		}
	}
}
