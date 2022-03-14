package telegram

import (
	"log"

	"github.com/zloyboy/gobot/internal/config"
	"github.com/zloyboy/gobot/internal/database"
	"github.com/zloyboy/gobot/internal/static"
	"github.com/zloyboy/gobot/internal/timeout"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot   *tgbotapi.BotAPI
	cfg   *config.Config
	dbase *database.Dbase
	stat  *static.Static
	tout  *timeout.Timeout
	uchan userChannel
}

func NewBot(bot *tgbotapi.BotAPI, db *database.Dbase, cfg *config.Config) *Bot {
	return &Bot{
		bot:   bot,
		cfg:   cfg,
		dbase: db,
		stat:  static.Stat(db),
		tout:  timeout.Make(),
		uchan: make(userChannel),
	}
}

func (b *Bot) Run() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	done := make(chan int64, 10)

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

	if b.cfg.Notify {
		chat := b.dbase.ReadChat()
		//chat := []int64{}
		log.Println("chat:", chat)
		for _, user := range chat {
			msg := tgbotapi.NewMessage(user, b.cfg.NotifyMsg)
			msg.ReplyMarkup = startKeyboard
			b.bot.Send(msg)
		}
	}

	for {
		select {
		case update := <-updates:
			if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
				continue
			}

			userID := update.SentFrom().ID
			if _, ok := b.uchan[userID]; !ok {
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

			if _, ok := b.uchan[userID]; ok {
				log.Printf("user %d data %s", userID, userData)
				if userData == "stop" || userData == "Stop" {
					b.uchan[userID].stop <- struct{}{}
				} else {
					b.uchan[userID].data <- userData
				}
			} else {
				chatID := update.FromChat().ID
				log.Printf("Start user %d", userID)
				if !b.dbase.ExistChat(chatID) {
					b.dbase.AddChat(chatID)
				}

				b.uchan[userID] = makeChannel()
				go RunSurvey(b, userID, chatID, b.uchan[userID], done, b.cfg.AnswerTout)
			}
		case userID := <-done:
			close(b.uchan[userID].stop)
			close(b.uchan[userID].data)
			delete(b.uchan, userID)
			log.Printf("Exit user %d", userID)
		}
	}
}
