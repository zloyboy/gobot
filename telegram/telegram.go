package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/zloyboy/gobot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var cntAll, cntYes, cntNo int
var ages = [6]string{"до 20", "20-29", "30-39", "40-49", "50-59", "60 ++"}
var ages_stat = [6][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}
var ageGroup = map[string]int{ages[0]: 15, ages[1]: 25, ages[2]: 35, ages[3]: 45, ages[4]: 55, ages[5]: 65}
var user_age = make(map[int64]int)
var user_timestamp = make(map[int64]int64)
var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
var start_msg = "Независимый подсчет статистики по COVID-19\nУкажите вашу возрастную группу:"

func inAges(age string) bool {
	switch age {
	case
		ages[0], ages[1], ages[2], ages[3], ages[4], ages[5]:
		return true
	}
	return false
}

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

var numericResKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "1"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "0"),
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

func (b *Bot) readStatFromDb() {
	cntAll = b.dbase.CountUsers()
	cntYes = b.dbase.CountRes()
	cntNo = cntAll - cntYes
	for i := 0; i < 6; i++ {
		ages_stat[i][0] = b.dbase.CountAge(i*10 + 15)
		if 0 < ages_stat[i][0] {
			ages_stat[i][1] = b.dbase.CountAgeRes(i*10 + 15)
		}
	}
}

func (b *Bot) makeStatFromDb() string {
	var perYes, perNo float32
	if cntAll == 0 {
		perYes, perNo = 0, 0
	} else {
		perYes = float32(cntYes) / float32(cntAll) * 100
		perNo = float32(cntNo) / float32(cntAll) * 100.0
	}
	perAge := [6]float32{0, 0, 0, 0, 0, 0}
	var outAge = ""
	for i := 0; i < 6; i++ {
		if 0 < ages_stat[i][0] {
			perAge[i] = float32(ages_stat[i][1]) / float32(ages_stat[i][0]) * 100
		}
		outAge += "\n" + ages[i] + " - " + fmt.Sprintf("%.2f", perAge[i]) + "% - " + strconv.Itoa(ages_stat[i][1]) + " из " + strconv.Itoa(ages_stat[i][0])
	}
	return "Независимая статистика по COVID-19\nОпрошено: " + strconv.Itoa(cntAll) +
		"\n" + fmt.Sprintf("%.2f", perYes) + "%" + " переболело: " + strconv.Itoa(cntYes) + " из " + strconv.Itoa(cntAll) +
		"\n" + fmt.Sprintf("%.2f", perNo) + "%" + " не болело: " + strconv.Itoa(cntNo) + " из " + strconv.Itoa(cntAll) +
		"\nЗаболеваемость по возрастным группам:" + outAge
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
				msg.Text = b.makeStatFromDb() + "\n---------------\nВы уже приняли участие в подсчете под именем " + uname + repeat_msg
				msg.ReplyMarkup = startKeyboard
			} else {
				// new user - send age question
				user_age[userID] = 0
				msg.Text = start_msg
				msg.ReplyMarkup = numericInlineKeyboard
			}
		} else {
			// user_age must have key=userID
			if uage, ok := user_age[userID]; ok {
				if uage == 0 {
					// get age - send res question
					age := userData
					if inAges(age) {
						user_age[userID] = ageGroup[age]
						msg.Text = "Вы переболели covid19?\n(по официальному мед.заключению)"
						msg.ReplyMarkup = numericResKeyboard
					} else {
						msg.Text = "Произошла ошибка: неверный возраст " + age
					}
				} else {
					// write poll results to DB
					outMsg := "Произошла ошибка: повторный ввод"
					res, _ := strconv.Atoi(userData)
					if b.dbase.NewId(userID) {
						log.Printf("insert id %d, name %s, age %d, res %d", userID, userName, user_age[userID], res)
						b.dbase.Insert(userID,
							time.Now().Local().Format("2006-01-02 15:04:05"),
							userName,
							user_age[userID],
							res)
						outMsg = b.makeStatFromDb()
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
