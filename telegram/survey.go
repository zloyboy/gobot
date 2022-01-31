package telegram

import (
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const survey_default = 0

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
}

func MakeUser() UserData {
	return UserData{country: "",
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
		monthVac:  nil}
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
	return &UserSession{b, survey_default, 0, 0, userID, chatID, userName, make(chan string), MakeUser()}
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

func (s *UserSession) startSurvey_00() bool {
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

func (s *UserSession) askCountry_00() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_country_msg
	msg.ReplyMarkup = countryInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getCountry_01(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	country := userData
	switch country {
	case Russia[1], Ukraine[1], Belarus[1], Kazakh[1]:
		s.userData.country = country
	default:
		msg.Text = error_msg + "страна " + country + " не участвует в опросе"
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askBirth_01() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_birth_msg
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getBirth_02(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	year, _ := strconv.Atoi(userData)
	if 1920 < year && year < 2020 {
		s.userData.birth = year
	} else {
		msg.Text = error_msg + userData + " не корректный год"
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askGender_02() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_gender_msg
	msg.ReplyMarkup = genderInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getGender_03(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	if userData == Male[1] || userData == Female[1] {
		s.userData.gender, _ = strconv.Atoi(userData)
	} else {
		msg.Text = error_msg + "не корректный пол"
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askEducation_03() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_education_msg
	msg.ReplyMarkup = educationInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getEducation_04(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	education := userData
	switch education {
	case School[1], College[1], University[1]:
		s.userData.education = education
	default:
		msg.Text = error_msg + education + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askOrigin_04() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_origin_msg
	msg.ReplyMarkup = originInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getOrigin_05(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	origin := userData
	switch origin {
	case Nature[1], Human[1], Unknown[1]:
		s.userData.origin = origin
	default:
		msg.Text = error_msg + origin + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askVaccine_05() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_vaccine_msg
	msg.ReplyMarkup = vaccineInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getVaccine_06(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	vaccine := userData
	switch vaccine {
	case Helpful[1], Useless[1], Dangerous[1], Unknown[1]:
		s.userData.vaccine = vaccine
	default:
		msg.Text = error_msg + vaccine + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askHaveIll_06() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_haveill_msg
	msg.ReplyMarkup = yesnoInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getHaveIll_07(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

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

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askCountIll_07() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_countill_msg
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getCountIll_08(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	countIll, _ := strconv.Atoi(userData)
	if 0 < countIll && countIll < 5 {
		s.userData.countIll = countIll
	} else {
		msg.Text = error_msg + userData + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askYearIll_09() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_yearill_msg + strconv.Itoa(s.count+1) + "й раз"
	msg.ReplyMarkup = yearillInlineKeyboard
	s.nextSubStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getYearIll_09(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	yearIll, _ := strconv.Atoi(userData)
	if 2020 <= yearIll && yearIll <= 2022 {
		s.userData.yearIll = append(s.userData.yearIll, yearIll)
	} else {
		msg.Text = error_msg + userData + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askMonthIll_09() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_monthill_msg + strconv.Itoa(s.count+1) + "й раз"
	msg.ReplyMarkup = monthillInlineKeyboard
	s.nextSubStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getMonthIll_09(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	monthIll, _ := strconv.Atoi(userData)
	if 1 <= monthIll && monthIll <= 12 {
		s.userData.monthIll = append(s.userData.monthIll, monthIll)
	} else {
		msg.Text = error_msg + userData + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askSignIll_09() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_signill_msg
	msg.ReplyMarkup = signillInlineKeyboard
	s.nextSubStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getSignIll_09(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	signIll := userData
	switch signIll {
	case Medic[1], Test[1], Symptom[1]:
		s.userData.signIll = append(s.userData.signIll, signIll)
	default:
		msg.Text = error_msg + userData + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askDegreeIll_09() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_degreeill_msg
	msg.ReplyMarkup = degreeillInlineKeyboard
	s.nextSubStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getDegreeIll_09(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

	degreeIll := userData
	switch degreeIll {
	case HospIvl[1], Hospital[1], HomeHard[1], HomeEasy[1], OnFoot[1], NoSymptom[1]:
		s.userData.degreeIll = append(s.userData.degreeIll, degreeIll)
	default:
		msg.Text = error_msg + degreeIll + error_ans
		ok = false
	}

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askHaveVac_09() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_havevac_msg
	msg.ReplyMarkup = yesnoInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getHaveVac_10(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var ok = true

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

	s.b.bot.Send(msg)
	return ok
}

func (s *UserSession) askCountVac_10() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_countvac_msg
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) writeResult(userData string) {
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
	s.b.utime.SetStamp(s.userID, time.Now().Unix())
	delete(user_session, s.userID)
}

func (s *UserSession) RunSurvey(ch chan string) {
	defer s.exit()
	if s.startSurvey_00() {
		s.askCountry_00()
	}
	for {
		data := <-ch
		switch s.state {
		case 1:
			if !s.getCountry_01(data) {
				s.abort()
				return
			}
			s.askBirth_01()
		case 2:
			if !s.getBirth_02(data) {
				s.abort()
				return
			}
			s.askGender_02()
		case 3:
			if !s.getGender_03(data) {
				s.abort()
				return
			}
			s.askEducation_03()
		case 4:
			if !s.getEducation_04(data) {
				s.abort()
				return
			}
			s.askOrigin_04()
		case 5:
			if !s.getOrigin_05(data) {
				s.abort()
				return
			}
			s.askVaccine_05()
		case 6:
			if !s.getVaccine_06(data) {
				s.abort()
				return
			}
			s.askHaveIll_06()
		case 7:
			if !s.getHaveIll_07(data) {
				s.abort()
				return
			}
			if 0 < s.userData.countIll {
				s.askCountIll_07()
			} else {
				s.nextStep()
				s.nextStep()
				s.askHaveVac_09()
			}
		case 8:
			if !s.getCountIll_08(data) {
				s.abort()
				return
			}
			s.count = 0
			s.resetSubStep()
			s.askYearIll_09()
			s.nextStep()
		case 9:
			switch s.subState {
			case 1:
				if !s.getYearIll_09(data) {
					s.abort()
					return
				}
				s.askMonthIll_09()
			case 2:
				if !s.getMonthIll_09(data) {
					s.abort()
					return
				}
				s.askSignIll_09()
			case 3:
				if !s.getSignIll_09(data) {
					s.abort()
					return
				}
				s.askDegreeIll_09()
			case 4:
				if !s.getDegreeIll_09(data) {
					s.abort()
					return
				}
				s.resetSubStep()
				s.count++
				if s.count < s.userData.countIll {
					s.askYearIll_09()
				} else {
					s.askHaveVac_09()
				}
			}
		case 10:
			if !s.getHaveVac_10(data) {
				s.abort()
				return
			}
			if 0 < s.userData.countVac {
				s.askCountVac_10()
			}
		default:
			s.writeResult(data)
			return
		}
	}
}
