package telegram

import (
	"log"
	"strconv"
	"time"

	"github.com/zloyboy/gobot/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const st_country = 0
const st_birth = 1
const st_gender = 2
const st_education = 3
const st_vacc_opin = 4
const st_orgn_opin = 5
const st_have_ill = st_orgn_opin + 1 //6

//const sst_year = user.Idx_year
const sst_month = user.Idx_month
const sst_sign = user.Idx_sign
const sst_degree = user.Idx_degree
const sst_kind = user.Idx_kind
const sst_effect = user.Idx_effect
const sst_next = sst_degree + 1 // 4

var user_session = make(map[int64]*UserSession)

type UserSession struct {
	b              *Bot
	state          int
	subState       int
	count          int
	userID, chatID int64
	userName       string
	userChan       chan string
	userData       user.UserData
	userIll        [4]int
	userVac        [4]int
}

func MakeSession(b *Bot, userID, chatID int64, userName string) *UserSession {
	return &UserSession{b, 0, 0, 0, userID, chatID, userName, make(chan string), user.MakeUser(), user.MakeSubUser(), user.MakeSubUser()}
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

func (s *UserSession) startSurvey() bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var res bool

	uname, err := s.b.dbase.CheckIdName(s.userID)
	if err == nil {
		// exist user - send statistic
		msg.Text = s.b.stat.MakeStatic() + "\nВы уже приняли участие в подсчете под именем " + uname + repeat_msg
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

func (s *UserSession) sendQuestion(text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(s.chatID, text)
	msg.ReplyMarkup = keyboard
	s.b.bot.Send(msg)
	s.nextStep()
}

func (s *UserSession) sendSubQuestion(text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(s.chatID, text)
	msg.ReplyMarkup = keyboard
	s.b.bot.Send(msg)
	s.nextSubStep()
}

func (s *UserSession) sendRequest() bool {
	switch s.state {
	case st_birth, st_gender, st_education, st_vacc_opin, st_orgn_opin, st_have_ill:
		s.sendQuestion(baseQuestion[s.state].ask, baseQuestion[s.state].key)
	case 7:
		if 0 < s.userData.CountIll {
			s.sendQuestion(ask_countill_msg, nil)
		} else {
			s.nextStep()
			s.nextStep()
			s.sendQuestion(ask_havevac_msg, yesnoInlineKeyboard)
		}
	case 8:
		s.count = 0
		s.resetSubStep()
		s.sendSubQuestion(illQuestion[s.subState].ask, illQuestion[s.subState].key)
		s.nextStep()
	case 9:
		switch s.subState {
		case sst_month, sst_sign, sst_degree:
			s.sendSubQuestion(illQuestion[s.subState].ask, illQuestion[s.subState].key)
		case 4:
			s.resetSubStep()
			s.count++
			s.userData.Ill = append(s.userData.Ill, s.userIll)
			if s.count < s.userData.CountIll {
				s.sendSubQuestion(illQuestion[s.subState].ask, yearInlineKeyboard)
			} else {
				s.sendQuestion(ask_havevac_msg, yesnoInlineKeyboard)
			}
		}
	case 10:
		if 0 < s.userData.CountVac {
			s.sendQuestion(ask_countvac_msg, nil)
		} else {
			return false
		}
	case 11:
		s.count = 0
		s.resetSubStep()
		s.sendSubQuestion(vacQuestion[s.subState].ask, vacQuestion[s.subState].key)
		s.nextStep()
	case 12:
		switch s.subState {
		case sst_month, sst_kind, sst_effect:
			s.sendSubQuestion(vacQuestion[s.subState].ask, vacQuestion[s.subState].key)
		case 4:
			s.resetSubStep()
			s.count++
			s.userData.Vac = append(s.userData.Vac, s.userVac)
			if s.count < s.userData.CountVac {
				s.sendSubQuestion(vacQuestion[s.subState].ask, vacQuestion[s.subState].key)
			} else {
				return false
			}
		}
	}

	return true
}

func (s *UserSession) getAnswer(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	switch s.state {
	case st_birth, st_gender, st_education, st_vacc_opin, st_orgn_opin, st_have_ill:
		val, _ := strconv.Atoi(userData)
		idx := s.state - 1
		if baseQuestion[idx].min <= val && val <= baseQuestion[idx].max {
			s.userData.Base[idx] = val
		} else {
			msg.Text = error_msg + error_ans
			ok = false
		}
	case 7:
		haveill := userData
		switch haveill {
		case Yes[1]:
			s.userData.CountIll = 1
		case No[1]:
			s.userData.CountIll = 0
		default:
			msg.Text = error_msg + haveill + error_ans
			ok = false
		}
	case 8:
		countIll, _ := strconv.Atoi(userData)
		if 0 < countIll && countIll < 5 {
			s.userData.CountIll = countIll
		} else {
			msg.Text = error_msg + userData + error_ans
			ok = false
		}
	case 9:
		switch s.subState {
		case sst_month, sst_sign, sst_degree, sst_next:
			val, _ := strconv.Atoi(userData)
			idx := s.subState - 1
			if illQuestion[idx].min <= val && val <= illQuestion[idx].max {
				s.userIll[idx] = val
			} else {
				msg.Text = error_msg + error_ans
				ok = false
			}
		}
	case 10:
		haveVac := userData
		switch haveVac {
		case Yes[1]:
			s.userData.CountVac = 1
		case No[1]:
			s.userData.CountVac = 0
		default:
			msg.Text = error_msg + haveVac + error_ans
			ok = false
		}
	case 11:
		countVac, _ := strconv.Atoi(userData)
		if 0 < countVac && countVac < 5 {
			s.userData.CountVac = countVac
		} else {
			msg.Text = error_msg + userData + error_ans
			ok = false
		}
	case 12:
		switch s.subState {
		case sst_month, sst_kind, sst_effect, sst_next:
			val, _ := strconv.Atoi(userData)
			idx := s.subState - 1
			if vacQuestion[idx].min <= val && val <= vacQuestion[idx].max {
				s.userVac[idx] = val
			} else {
				msg.Text = error_msg + error_ans
				ok = false
			}
		}
	}

	s.b.bot.Send(msg)
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
		s.userData.Base[st_country],
		s.userData.Base[st_birth],
		s.userData.Base[st_gender],
		s.userData.Base[st_education],
		s.userData.Base[st_vacc_opin],
		s.userData.Base[st_orgn_opin],
		s.userData.CountIll,
		s.userData.Ill,
		s.userData.CountVac,
		s.userData.Vac)
	s.b.stat.RefreshStatic(s.userData.CountIll)
	msg.Text = s.b.stat.MakeStatic() + repeat_msg
	msg.ReplyMarkup = startKeyboard

	s.b.bot.Send(msg)
}

func (s *UserSession) abort() {
	msg := tgbotapi.NewMessage(s.chatID, abort_msg)
	msg.ReplyMarkup = startKeyboard
	s.b.bot.Send(msg)
}

func (s *UserSession) exit() {
	log.Printf("Exit user %d", s.userID)
	s.b.utime.SetStamp(s.userID, time.Now().Unix())
	delete(user_session, s.userID)
}

func (s *UserSession) RunSurvey(ch chan string) {
	defer s.exit()
	if s.startSurvey() {
		s.sendQuestion(baseQuestion[st_country].ask, baseQuestion[st_country].key)
		for {
			data := <-ch
			if !s.getAnswer(data) {
				s.abort()
				return
			}
			if !s.sendRequest() {
				break
			}
		}
		s.writeResult()
	}
}
