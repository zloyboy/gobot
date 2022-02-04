package telegram

import (
	"log"
	"strconv"
	"time"

	"github.com/zloyboy/gobot/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	userIll        user.UserIll
	userVac        user.UserVac
}

func MakeSession(b *Bot, userID, chatID int64, userName string) *UserSession {
	return &UserSession{b, 0, 0, 0, userID, chatID, userName, make(chan string), user.MakeUser(), user.MakeIll(), user.MakeVac()}
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
	case 1:
		s.sendQuestion(ask_birth_msg, nil)
	case 2:
		s.sendQuestion(ask_gender_msg, genderInlineKeyboard)
	case 3:
		s.sendQuestion(ask_education_msg, educationInlineKeyboard)
	case 4:
		s.sendQuestion(ask_vaccine_msg, vaccineInlineKeyboard)
	case 5:
		s.sendQuestion(ask_origin_msg, originInlineKeyboard)
	case 6:
		s.sendQuestion(ask_haveill_msg, yesnoInlineKeyboard)
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
		s.sendSubQuestion(ask_yearill_msg+nTimes(s.count+1, s.userData.CountIll), yearInlineKeyboard)
		s.nextStep()
	case 9:
		switch s.subState {
		case 1:
			s.sendSubQuestion(ask_monthill_msg+nTimes(s.count+1, s.userData.CountIll), monthInlineKeyboard)
		case 2:
			s.sendSubQuestion(ask_signill_msg, signillInlineKeyboard)
		case 3:
			s.sendSubQuestion(ask_degreeill_msg, degreeillInlineKeyboard)
		case 4:
			s.resetSubStep()
			s.count++
			s.userData.Ill = append(s.userData.Ill, s.userIll)
			if s.count < s.userData.CountIll {
				s.sendSubQuestion(ask_yearill_msg+nTimes(s.count+1, s.userData.CountIll), yearInlineKeyboard)
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
		s.sendSubQuestion(ask_yearvac_msg+nTimes(s.count+1, s.userData.CountVac), yearInlineKeyboard)
		s.nextStep()
	case 12:
		switch s.subState {
		case 1:
			s.sendSubQuestion(ask_monthvac_msg+nTimes(s.count+1, s.userData.CountVac), monthInlineKeyboard)
		case 2:
			s.sendSubQuestion(ask_kindvac_msg, kindvacInlineKeyboard)
		case 3:
			s.sendSubQuestion(ask_effectvac_msg, effectvacInlineKeyboard)
		case 4:
			s.resetSubStep()
			s.count++
			s.userData.Vac = append(s.userData.Vac, s.userVac)
			if s.count < s.userData.CountVac {
				s.sendSubQuestion(ask_yearvac_msg+nTimes(s.count+1, s.userData.CountVac), yearInlineKeyboard)
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
	case 1:
		country, _ := strconv.Atoi(userData)
		if country < 4 {
			s.userData.Country = country
		} else {
			msg.Text = error_msg + error_ans
			ok = false
		}
	case 2:
		year, _ := strconv.Atoi(userData)
		if 1920 < year && year < 2020 {
			s.userData.Birth = year
		} else {
			msg.Text = error_msg + userData + " не корректный год"
			ok = false
		}
	case 3:
		if userData == Male[1] || userData == Female[1] {
			s.userData.Gender, _ = strconv.Atoi(userData)
		} else {
			msg.Text = error_msg + "не корректный пол"
			ok = false
		}
	case 4:
		education, _ := strconv.Atoi(userData)
		if education < 3 {
			s.userData.Education = education
		} else {
			msg.Text = error_msg + error_ans
			ok = false
		}
	case 5:
		vaccine, _ := strconv.Atoi(userData)
		if vaccine < 3 {
			s.userData.Vaccine = vaccine
		} else {
			msg.Text = error_msg + error_ans
			ok = false
		}
	case 6:
		origin, _ := strconv.Atoi(userData)
		if origin < 2 {
			s.userData.Origin = origin
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
		case 1:
			yearIll, _ := strconv.Atoi(userData)
			if 2020 <= yearIll && yearIll <= 2022 {
				s.userIll.Year = yearIll
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 2:
			monthIll, _ := strconv.Atoi(userData)
			if 1 <= monthIll && monthIll <= 12 {
				s.userIll.Month = monthIll
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 3:
			signIll, _ := strconv.Atoi(userData)
			if signIll < 3 {
				s.userIll.Sign = signIll
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 4:
			degreeIll, _ := strconv.Atoi(userData)
			if degreeIll < 6 {
				s.userIll.Degree = degreeIll
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
		case 1:
			yearVac, _ := strconv.Atoi(userData)
			if 2020 <= yearVac && yearVac <= 2022 {
				s.userVac.Year = yearVac
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 2:
			monthVac, _ := strconv.Atoi(userData)
			if 1 <= monthVac && monthVac <= 12 {
				s.userVac.Month = monthVac
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 3:
			kindVac := userData
			switch kindVac {
			case SputnikV[1], SputnikL[1], EpiVac[1], Kovivak[1]:
				s.userVac.Kind = kindVac
			default:
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 4:
			effectVac := userData
			switch effectVac {
			case HardEffect[1], MediumEffect[1], EasyEffect[1]:
				s.userVac.Effect = effectVac
			default:
				msg.Text = error_msg + userData + error_ans
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
		s.userID, s.userName, s.userData.Country, s.userData.Birth, s.userData.Gender, s.userData.Education, s.userData.Vaccine,
		s.userData.Origin, s.userData.CountIll, s.userData.CountVac)
	s.b.dbase.Insert(s.userID,
		time.Now().Local().Format("2006-01-02 15:04:05"),
		s.userName,
		s.userData.Country,
		s.userData.Birth,
		s.userData.Gender,
		s.userData.Education,
		s.userData.Vaccine,
		s.userData.Origin,
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
		s.sendQuestion(ask_country_msg, countryInlineKeyboard)
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
