package telegram

import (
	"log"
	"time"

	"github.com/zloyboy/gobot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var user_age_idx = make(map[int64]int)
var user_timestamp = make(map[int64]int64)
var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
var start_msg = "Независимый подсчет статистики по COVID-19\nУкажите вашу возрастную группу:"

type Bot struct {
	bot   *tgbotapi.BotAPI
	dbase *database.Dbase
	stat  *Static
}

func NewBot(bot *tgbotapi.BotAPI, db *database.Dbase) *Bot {
	return &Bot{bot: bot, dbase: db, stat: Stat()}
}

func (b *Bot) userTimeout(userID int64) bool {
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
		userName := update.SentFrom().FirstName
		chatID := update.FromChat().ID
		var userData string
		if update.Message != nil {
			userData = update.Message.Text
		} else if update.CallbackQuery != nil {
			userData = update.CallbackQuery.Data
		}
		log.Printf("user %d, data %s", userID, userData)

		if update.Message != nil {
			if b.userTimeout(userID) {
				continue
			}
			b.askAge_01(userID, chatID)
		} else {
			// user_age_idx must have key=userID
			if uage, ok := user_age_idx[userID]; ok {
				if uage == -1 {
					// get age - send res question
					b.askIll_02(userID, chatID, userData)
				} else {
					// write poll results to DB
					b.writeResult_03(userID, chatID, userData, userName)
				}
			} else {
				b.internalError(chatID)
			}
		}
	}

	return nil
}
