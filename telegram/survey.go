package telegram

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var user_session = make(map[int64]*UserSession)

type UserData struct {
	country   string
	birth     int
	gender    int
	education string
	origin    string
	vaccine   string
	countIll  int
	yearIll   []int
	monthIll  []int
	signIll   []string
	degreeIll []string
	countVac  int
	yearVac   []int
	monthVac  []int
	kindVac   []string
	effectVac []string
}

func MakeUser() UserData {
	return UserData{
		country:   "",
		birth:     -1,
		gender:    -1,
		education: "",
		origin:    "",
		vaccine:   "",
		countIll:  0,
		yearIll:   nil,
		monthIll:  nil,
		signIll:   nil,
		degreeIll: nil,
		countVac:  0,
		yearVac:   nil,
		monthVac:  nil,
		kindVac:   nil,
		effectVac: nil}
}

type UserSession struct {
	b              *Bot
	state          int
	subState       int
	count          int
	userID, chatID int64
	userName       string
	userChan       chan string
	userData       UserData
}

func MakeSession(b *Bot, userID, chatID int64, userName string) *UserSession {
	return &UserSession{b, 0, 0, 0, userID, chatID, userName, make(chan string), MakeUser()}
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
		msg.Text = s.b.stat.MakeStatic() + "\n---------------\nВы уже приняли участие в подсчете под именем " + uname + repeat_msg
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
		s.sendQuestion(ask_origin_msg, originInlineKeyboard)
	case 5:
		s.sendQuestion(ask_vaccine_msg, vaccineInlineKeyboard)
	case 6:
		s.sendQuestion(ask_haveill_msg, yesnoInlineKeyboard)
	case 7:
		if 0 < s.userData.countIll {
			s.sendQuestion(ask_countill_msg, nil)
		} else {
			s.nextStep()
			s.nextStep()
			s.sendQuestion(ask_havevac_msg, yesnoInlineKeyboard)
		}
	case 8:
		s.count = 0
		s.resetSubStep()
		s.sendSubQuestion(ask_yearill_msg+strconv.Itoa(s.count+1)+"й раз", yearInlineKeyboard)
		s.nextStep()
	case 9:
		switch s.subState {
		case 1:
			s.sendSubQuestion(ask_monthill_msg+strconv.Itoa(s.count+1)+"й раз", monthInlineKeyboard)
		case 2:
			s.sendSubQuestion(ask_signill_msg, signillInlineKeyboard)
		case 3:
			s.sendSubQuestion(ask_degreeill_msg, degreeillInlineKeyboard)
		case 4:
			s.resetSubStep()
			s.count++
			if s.count < s.userData.countIll {
				s.sendSubQuestion(ask_yearill_msg+strconv.Itoa(s.count+1)+"й раз", yearInlineKeyboard)
			} else {
				s.sendQuestion(ask_havevac_msg, yesnoInlineKeyboard)
			}
		}
	case 10:
		if 0 < s.userData.countVac {
			s.sendQuestion(ask_countvac_msg, nil)
		} else {
			return false
		}
	case 11:
		s.count = 0
		s.resetSubStep()
		s.sendSubQuestion(ask_yearvac_msg+strconv.Itoa(s.count+1)+"й раз", yearInlineKeyboard)
		s.nextStep()
	case 12:
		switch s.subState {
		case 1:
			s.sendSubQuestion(ask_monthvac_msg+strconv.Itoa(s.count+1)+"й раз", monthInlineKeyboard)
		case 2:
			s.sendSubQuestion(ask_kindvac_msg, kindvacInlineKeyboard)
		case 3:
			s.sendSubQuestion(ask_effectvac_msg, effectvacInlineKeyboard)
		case 4:
			s.resetSubStep()
			s.count++
			if s.count < s.userData.countVac {
				s.sendSubQuestion(ask_yearvac_msg+strconv.Itoa(s.count+1)+"й раз", yearInlineKeyboard)
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
		country := userData
		switch country {
		case Russia[1], Ukraine[1], Belarus[1], Kazakh[1]:
			s.userData.country = country
		default:
			msg.Text = error_msg + "страна " + country + " не участвует в опросе"
			ok = false
		}
	case 2:
		year, _ := strconv.Atoi(userData)
		if 1920 < year && year < 2020 {
			s.userData.birth = year
		} else {
			msg.Text = error_msg + userData + " не корректный год"
			ok = false
		}
	case 3:
		if userData == Male[1] || userData == Female[1] {
			s.userData.gender, _ = strconv.Atoi(userData)
		} else {
			msg.Text = error_msg + "не корректный пол"
			ok = false
		}
	case 4:
		education := userData
		switch education {
		case School[1], College[1], University[1]:
			s.userData.education = education
		default:
			msg.Text = error_msg + education + error_ans
			ok = false
		}
	case 5:
		origin := userData
		switch origin {
		case Nature[1], Human[1], Unknown[1]:
			s.userData.origin = origin
		default:
			msg.Text = error_msg + origin + error_ans
			ok = false
		}
	case 6:
		vaccine := userData
		switch vaccine {
		case Helpful[1], Useless[1], Dangerous[1], Unknown[1]:
			s.userData.vaccine = vaccine
		default:
			msg.Text = error_msg + vaccine + error_ans
			ok = false
		}
	case 7:
		haveill := userData
		switch haveill {
		case Yes[1]:
			s.userData.countIll = 1
		case No[1]:
			s.userData.countIll = 0
		default:
			msg.Text = error_msg + haveill + error_ans
			ok = false
		}
	case 8:
		countIll, _ := strconv.Atoi(userData)
		if 0 < countIll && countIll < 5 {
			s.userData.countIll = countIll
		} else {
			msg.Text = error_msg + userData + error_ans
			ok = false
		}
	case 9:
		switch s.subState {
		case 1:
			yearIll, _ := strconv.Atoi(userData)
			if 2020 <= yearIll && yearIll <= 2022 {
				s.userData.yearIll = append(s.userData.yearIll, yearIll)
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 2:
			monthIll, _ := strconv.Atoi(userData)
			if 1 <= monthIll && monthIll <= 12 {
				s.userData.monthIll = append(s.userData.monthIll, monthIll)
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 3:
			signIll := userData
			switch signIll {
			case Medic[1], Test[1], Symptom[1]:
				s.userData.signIll = append(s.userData.signIll, signIll)
			default:
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 4:
			degreeIll := userData
			switch degreeIll {
			case HospIvl[1], Hospital[1], HomeHard[1], HomeEasy[1], OnFoot[1], NoSymptom[1]:
				s.userData.degreeIll = append(s.userData.degreeIll, degreeIll)
			default:
				msg.Text = error_msg + degreeIll + error_ans
				ok = false
			}
		}
	case 10:
		haveVac := userData
		switch haveVac {
		case Yes[1]:
			s.userData.countVac = 1
		case No[1]:
			s.userData.countVac = 0
		default:
			msg.Text = error_msg + haveVac + error_ans
			ok = false
		}
	case 11:
		countVac, _ := strconv.Atoi(userData)
		if 0 < countVac && countVac < 5 {
			s.userData.countVac = countVac
		} else {
			msg.Text = error_msg + userData + error_ans
			ok = false
		}
	case 12:
		switch s.subState {
		case 1:
			yearVac, _ := strconv.Atoi(userData)
			if 2020 <= yearVac && yearVac <= 2022 {
				s.userData.yearVac = append(s.userData.yearVac, yearVac)
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 2:
			monthVac, _ := strconv.Atoi(userData)
			if 1 <= monthVac && monthVac <= 12 {
				s.userData.monthVac = append(s.userData.monthVac, monthVac)
			} else {
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 3:
			kindVac := userData
			switch kindVac {
			case SputnikV[1], SputnikL[1], EpiVac[1], Kovivak[1]:
				s.userData.kindVac = append(s.userData.kindVac, kindVac)
			default:
				msg.Text = error_msg + userData + error_ans
				ok = false
			}
		case 4:
			effectVac := userData
			switch effectVac {
			case HardEffect[1], MediumEffect[1], EasyEffect[1]:
				s.userData.effectVac = append(s.userData.effectVac, effectVac)
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

	//ill, _ := strconv.Atoi(userData)
	if s.b.dbase.NewId(s.userID) {
		/*aver_age := age_mid[s.age_idx]
		log.Printf("insert id %d, name %s, age %d, ill %d", s.userID, s.userName, aver_age, ill)
		s.b.dbase.Insert(s.userID,
			time.Now().Local().Format("2006-01-02 15:04:05"),
			s.userName,
			aver_age,
			ill)
		s.b.stat.RefreshStatic(s.age_idx, ill)*/
		msg.Text = s.b.stat.MakeStatic()
	} else {
		msg.Text = "Произошла ошибка: повторный ввод"
	}

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
	}
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
