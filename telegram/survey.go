package telegram

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) askAge_01(userID, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "")

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

	b.bot.Send(msg)
}

func (b *Bot) askIll_02(userID, chatID int64, userData string) {
	msg := tgbotapi.NewMessage(chatID, "")

	age := userData
	if inAges(age) {
		user_age_idx[userID] = age_idx[age]
		msg.Text = "Вы переболели covid19?\n(по официальному мед.заключению)"
		msg.ReplyMarkup = numericResKeyboard
	} else {
		msg.Text = "Произошла ошибка: неверный возраст " + age
	}

	b.bot.Send(msg)
}

func (b *Bot) writeResult_03(userID, chatID int64, userData, userName string) {
	msg := tgbotapi.NewMessage(chatID, "")

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
		msg.Text = b.stat.MakeStatic()
	} else {
		msg.Text = "Произошла ошибка: повторный ввод"
	}

	b.bot.Send(msg)
}

func (b *Bot) internalError(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Произошла ошибка сервера")
	b.bot.Send(msg)
}
