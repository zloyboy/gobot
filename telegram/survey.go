package telegram

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
var start_msg = "Независимый подсчет статистики по COVID-19\nУкажите вашу возрастную группу:"

var user_session = make(map[int64]*UserSession)
var user_timestamp = make(map[int64]int64)

type UserSession struct {
	b              *Bot
	userID, chatID int64
	userName       string
	age_idx        int
	userChan       chan string
}

func (s *UserSession) askAge_01() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	uname, err := s.b.dbase.CheckIdName(s.userID)
	if err == nil {
		// exist user - send statistic
		msg.Text = s.b.stat.MakeStatic() + "\n---------------\nВы уже приняли участие в подсчете под именем " + uname + repeat_msg
		msg.ReplyMarkup = startKeyboard
	} else {
		// new user - send age question
		s.age_idx = -1
		msg.Text = start_msg
		msg.ReplyMarkup = numericInlineKeyboard
		user_timestamp[s.userID] = 0
	}

	s.b.bot.Send(msg)
}

func (s *UserSession) askIll_02(userData string) {
	msg := tgbotapi.NewMessage(s.chatID, "")

	age := userData
	if inAges(age) {
		s.age_idx = age_idx[age]
		msg.Text = "Вы переболели covid19?\n(по официальному мед.заключению)"
		msg.ReplyMarkup = numericResKeyboard
	} else {
		msg.Text = "Произошла ошибка: неверный возраст " + age
	}

	s.b.bot.Send(msg)
	user_timestamp[s.userID] = 0
}

func (s *UserSession) writeResult_03(userData string) {
	msg := tgbotapi.NewMessage(s.chatID, "")

	ill, _ := strconv.Atoi(userData)
	if s.b.dbase.NewId(s.userID) {
		aver_age := age_mid[s.age_idx]
		log.Printf("insert id %d, name %s, age %d, ill %d", s.userID, s.userName, aver_age, ill)
		s.b.dbase.Insert(s.userID,
			time.Now().Local().Format("2006-01-02 15:04:05"),
			s.userName,
			aver_age,
			ill)
		s.b.stat.RefreshStatic(s.age_idx, ill)
		msg.Text = s.b.stat.MakeStatic()
	} else {
		msg.Text = "Произошла ошибка: повторный ввод"
	}

	s.b.bot.Send(msg)
}

func (s *UserSession) exit() {
	user_timestamp[s.userID] = time.Now().Unix()
	delete(user_session, s.userID)
}

func (s *UserSession) RunSurvey(ch chan string) {
	defer s.exit()
	s.askAge_01()
	for {
		data := <-ch
		switch s.age_idx {
		case -2:
			s.askAge_01()
			return
		case -1:
			s.askIll_02(data)
		default:
			s.writeResult_03(data)
			return
		}
	}
}
