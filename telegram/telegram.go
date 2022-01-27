package telegram

import (
	"log"
	"strconv"
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

	b.readStatFromDb()

	var chatID int64
	var userID int64
	var userData string
	var userName string

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore non-Message updates
			continue
		}

		userID = update.SentFrom().ID
		userName = update.SentFrom().FirstName
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
			uname, err := b.dbase.CheckIdName(userID)
			if err == nil {
				// exist user - send statistic
				msg.Text = b.stat.MakeStatic() + "\n---------------\nВы уже приняли участие в подсчете под именем " + uname + repeat_msg
				msg.ReplyMarkup = startKeyboard
			} else {
				// new user - send age question
				user_age_idx[userID] = -1
				msg.Text = start_msg
				msg.ReplyMarkup = numericInlineKeyboard
			}
		} else {
			// user_age_idx must have key=userID
			if uage, ok := user_age_idx[userID]; ok {
				if uage == -1 {
					// get age - send res question
					age := userData
					if inAges(age) {
						user_age_idx[userID] = age_idx[age]
						msg.Text = "Вы переболели covid19?\n(по официальному мед.заключению)"
						msg.ReplyMarkup = numericResKeyboard
					} else {
						msg.Text = "Произошла ошибка: неверный возраст " + age
					}
				} else {
					// write poll results to DB
					outMsg := "Произошла ошибка: повторный ввод"
					ill, _ := strconv.Atoi(userData)
					if b.dbase.NewId(userID) {
						aver_age := age_mid[user_age_idx[userID]]
						log.Printf("insert id %d, name %s, age %d, ill %d", userID, userName, aver_age, ill)
						b.dbase.Insert(userID,
							time.Now().Local().Format("2006-01-02 15:04:05"),
							userName,
							aver_age,
							ill)
						b.stat.RefreshStatic(user_age_idx[userID], ill)
						outMsg = b.stat.MakeStatic()
					}
					msg.Text = outMsg
				}
			} else {
				msg.Text = "Произошла ошибка сервера"
			}
		}

		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

	return nil
}
