package telegram

import (
	"log"

	"github.com/zloyboy/gobot/internal/database"
	"github.com/zloyboy/gobot/internal/static"
	"github.com/zloyboy/gobot/internal/timeout"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot   *tgbotapi.BotAPI
	dbase *database.Dbase
	stat  *static.Static
	tout  *timeout.Timeout
}

func NewBot(bot *tgbotapi.BotAPI, db *database.Dbase) *Bot {
	return &Bot{bot: bot, dbase: db, stat: static.Stat(db), tout: timeout.Make()}
}

func (b *Bot) Run() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	if !b.readYearFromDb() ||
		!b.readMonthFromDb() ||
		!b.readCountryFromDb() ||
		!b.readEducationFromDb() ||
		!b.readVaccOpinionFromDb() ||
		!b.readOrgnOpinionFromDb() ||
		!b.readIllnessSignFromDb() ||
		!b.readIllnessDegreeFromDb() ||
		!b.readVaccineKindFromDb() ||
		!b.readVaccineEffectFromDb() {
		return
	}
	b.setKeyboards()
	b.stat.ReadStatFromDb()

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
			continue
		}

		userID := update.SentFrom().ID
		if _, ok := user_session[userID]; !ok {
			if b.tout.Exist(userID) {
				continue
			}
		}

		var userData string
		if update.Message != nil {
			userData = update.Message.Text
		} else if update.CallbackQuery != nil {
			userData = update.CallbackQuery.Data
		} else {
			continue
		}

		if _, ok := user_session[userID]; ok {
			//log.Printf("user %d data %s", userID, userData)
			if userData == "stop" || userData == "Stop" {
				close(user_session[userID].userStop)
			} else {
				user_session[userID].userChan <- userData
			}
		} else {
			chatID := update.FromChat().ID
			log.Printf("Start user %d", userID)

			user_session[userID] = MakeSession(b, userID, chatID)
			go user_session[userID].RunSurvey(
				user_session[userID].userChan,
				user_session[userID].userStop)
		}
	}
}
