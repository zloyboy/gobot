package telegram

import (
	"log"
	"time"

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

func (b *Bot) userTimeout(id int64) bool {
	curr_time := time.Now().Unix()
	if tstamp, ok := user_timestamp[id]; ok && 0 < tstamp {
		log.Printf("Timestamp %d, pass %d", tstamp, (curr_time - tstamp))
		if (curr_time - tstamp) < 10 {
			return true
		}
	}
	user_timestamp[id] = curr_time
	return false
}

func (b *Bot) Run() {
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
		if b.userTimeout(userID) {
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

			user_session[userID] = &UserSession{b, userID, chatID, userName, -2, make(chan string)}
			go user_session[userID].RunSurvey(user_session[userID].userChan)
		}
	}
}
