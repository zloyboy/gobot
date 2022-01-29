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
}

func NewBot(bot *tgbotapi.BotAPI, db *database.Dbase) *Bot {
	return &Bot{bot: bot, dbase: db, stat: Stat()}
}

func (b *Bot) Run() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	b.readStatFromDb()

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
			continue
		}

		userID := update.SentFrom().ID
		var userData string
		if update.Message != nil {
			userData = update.Message.Text
		} else if update.CallbackQuery != nil {
			userData = update.CallbackQuery.Data
		}

		if uses, ok := user_session[userID]; ok {
			log.Printf("doing user %d, data %s, name %s", userID, userData, uses.userName)
		} else {
			chatID := update.FromChat().ID
			userName := update.SentFrom().FirstName
			user_session[userID] = &UserSession{b, userID, chatID, 0, userName, -2}
			log.Printf("Start user %d, data %s", userID, userData)
		}

		if update.Message != nil {
			if user_session[userID].UserTimeout() {
				continue
			}
			user_session[userID].askAge_01()
		} else {
			if user_session[userID].GetAgeIdx() == -1 {
				// get age - send res question
				user_session[userID].askIll_02(userData)
			} else {
				// write poll results to DB
				user_session[userID].writeResult_03(userData)
				delete(user_session, userID)
			}
		}
	}

	return nil
}
