package telegram

import (
	"log"
	"strconv"
	"time"

	"github.com/zloyboy/gobot/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const st_country = user.Idx_country
const st_birth = user.Idx_birth
const st_gender = user.Idx_gender
const st_education = user.Idx_education
const st_vacc_opin = user.Idx_vacc_opin
const st_orgn_opin = user.Idx_orgn_opin
const st_have_ill = 6
const st_get_have_ill = 7
const st_get_count_ill = 8
const st_illness = 9
const st_get_have_vac = 10
const st_get_count_vac = 11
const st_vaccination = 12

const sst_year = user.Idx_year
const sst_month = user.Idx_month
const sst_sign = user.Idx_sign
const sst_degree = user.Idx_degree
const sst_kind = user.Idx_kind
const sst_effect = user.Idx_effect
const sst_next = 4

var user_session = make(map[int64]*UserSession)

type UserSession struct {
	b              *Bot
	state          int
	subState       int
	count          int
	userID, chatID int64
	userName       string
	userChan       chan string
	userStop       chan struct{}
	userData       user.UserData
	userSub        [4]int
}

func MakeSession(b *Bot, userID, chatID int64, userName string) *UserSession {
	return &UserSession{b,
		0, 0, 0,
		userID, chatID, userName,
		make(chan string, 10),
		make(chan struct{}),
		user.MakeUser(),
		user.MakeSubUser()}
}

func (s *UserSession) startSurvey() bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var res bool

	uname, err := s.b.dbase.CheckIdName(s.userID)
	if err == nil {
		// exist user - send statistic
		msg.Text = "\nВы уже приняли участие в подсчете под именем " + uname +
			"\n--------------------\n" + s.b.stat.MakeStatic() + repeat_msg
		msg.ReplyMarkup = startKeyboard
		res = false
	} else {
		// new user - start survey
		msg.Text = start_msg
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		res = true
	}

	s.b.bot.Send(msg)
	return res
}

func (s *UserSession) nextStep() {
	s.state++
	s.b.utime.ResetStamp(s.userID)
}

func (s *UserSession) resetSubStep() {
	s.subState = 0
	s.b.utime.ResetStamp(s.userID)
}

func (s *UserSession) nextSubStep() {
	s.subState++
	s.b.utime.ResetStamp(s.userID)
}

func (s *UserSession) sendQuestion(text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(s.chatID, text)
	msg.ReplyMarkup = keyboard
	s.b.bot.Send(msg)
	s.nextStep()
}

func (s *UserSession) sendSubQuestion(subQuestion *SubQuestion, opt ...int) {
	text := subQuestion.ask
	if subQuestion.many && 0 < len(opt) && 1 < opt[0] {
		text = subQuestion.ask + strconv.Itoa(s.count+1) + "й раз"
	}
	msg := tgbotapi.NewMessage(s.chatID, text)
	msg.ReplyMarkup = subQuestion.key
	s.b.bot.Send(msg)
	s.nextSubStep()
}

func (s *UserSession) sendRequest() bool {
	switch s.state {
	case st_birth, st_gender, st_education, st_vacc_opin, st_orgn_opin, st_have_ill:
		s.sendQuestion(baseQuestion[s.state].ask, baseQuestion[s.state].key)
	case st_get_have_ill:
		if 0 < s.userData.CountIll {
			s.sendQuestion(ask_countill_msg, numInlineKeyboard)
		} else {
			s.nextStep()
			s.nextStep()
			s.sendQuestion(ask_havevac_msg, yesnoInlineKeyboard)
		}
	case st_get_count_ill:
		s.count = 0
		s.resetSubStep()
		s.sendSubQuestion(&illQuestion[s.subState], s.userData.CountIll)
		s.nextStep()
	case st_illness:
		switch s.subState {
		case sst_month, sst_sign, sst_degree:
			s.sendSubQuestion(&illQuestion[s.subState], s.userData.CountIll)
		case sst_next:
			s.resetSubStep()
			s.count++
			s.userData.Ill = append(s.userData.Ill, s.userSub)
			if s.count < s.userData.CountIll {
				s.sendSubQuestion(&illQuestion[s.subState], s.userData.CountIll)
			} else {
				s.sendQuestion(ask_havevac_msg, yesnoInlineKeyboard)
			}
		}
	case st_get_have_vac:
		if 0 < s.userData.CountVac {
			s.sendQuestion(ask_countvac_msg, numInlineKeyboard)
		} else {
			return false
		}
	case st_get_count_vac:
		s.count = 0
		s.resetSubStep()
		s.sendSubQuestion(&vacQuestion[s.subState], s.userData.CountVac)
		s.nextStep()
	case st_vaccination:
		switch s.subState {
		case sst_month, sst_kind, sst_effect:
			s.sendSubQuestion(&vacQuestion[s.subState], s.userData.CountVac)
		case sst_next:
			s.resetSubStep()
			s.count++
			s.userData.Vac = append(s.userData.Vac, s.userSub)
			if s.count < s.userData.CountVac {
				s.sendSubQuestion(&vacQuestion[s.subState], s.userData.CountVac)
			} else {
				return false
			}
		}
	}

	return true
}

func (s *UserSession) getAnswer(userData string) bool {
	var ok = true

	switch s.state {
	case st_birth, st_gender, st_education, st_vacc_opin, st_orgn_opin, st_have_ill:
		val, err := strconv.Atoi(userData)
		idx := s.state - 1
		if err != nil || val < baseQuestion[idx].min || baseQuestion[idx].max < val {
			ok = false
		} else {
			s.userData.Base[idx] = val
		}
	case st_get_have_ill, st_get_have_vac:
		val, err := strconv.Atoi(userData)
		if err != nil || val != 0 && val != 1 {
			ok = false
		} else {
			if s.state == st_get_have_ill {
				s.userData.CountIll = val
			} else {
				s.userData.CountVac = val
			}
		}
	case st_get_count_ill, st_get_count_vac:
		val, err := strconv.Atoi(userData)
		if err != nil || val < 0 || 3 < val {
			ok = false
		} else {
			if s.state == st_get_count_ill {
				s.userData.CountIll = val
			} else {
				s.userData.CountVac = val
			}
		}
	case st_illness, st_vaccination:
		val, err := strconv.Atoi(userData)
		idx := s.subState - 1
		subQuestion := illQuestion
		if s.state == st_vaccination {
			subQuestion = vacQuestion
		}
		if err != nil || val < subQuestion[idx].min || subQuestion[idx].max < val {
			ok = false
		} else {
			s.userSub[idx] = val
		}
	}

	return ok
}

func (s *UserSession) writeResult() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	log.Printf("insert id %d, name %s, country %d, birth %d, gender %d, education %d, vaccine %d, origin %d, countIll %d, countVac %d",
		s.userID, s.userName, s.userData.Base[st_country], s.userData.Base[st_birth], s.userData.Base[st_gender], s.userData.Base[st_education],
		s.userData.Base[st_vacc_opin], s.userData.Base[st_orgn_opin], s.userData.CountIll, s.userData.CountVac)
	s.b.dbase.Insert(s.userID,
		time.Now().Local().Format("2006-01-02 15:04:05"),
		s.userName,
		s.userData)
	s.b.stat.RefreshStatic(s.userData)
	msg.Text = s.b.stat.MakeStatic() + repeat_msg
	msg.ReplyMarkup = startKeyboard

	s.b.bot.Send(msg)
}

func (s *UserSession) abort(reason string) {
	msg := tgbotapi.NewMessage(s.chatID, reason+again_msg)
	msg.ReplyMarkup = startKeyboard
	s.b.bot.Send(msg)
}

func (s *UserSession) exit() {
	log.Printf("Exit user %d", s.userID)
	s.b.utime.SetStamp(s.userID, time.Now().Unix())
	delete(user_session, s.userID)
}

func (s *UserSession) RunSurvey(ch chan string, quit chan struct{}) {
	defer s.exit()
	run := 1
	if s.startSurvey() {
		s.sendQuestion(baseQuestion[st_country].ask, baseQuestion[st_country].key)
		for run != 0 {
			select {
			case <-quit:
				s.abort(stop_msg)
				return
			case data := <-ch:
				if !s.getAnswer(data) {
					s.abort(error_msg)
					return
				}
				if !s.sendRequest() {
					run = 0
				}
			}
		}
		s.writeResult()
	}
}
