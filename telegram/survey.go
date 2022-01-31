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
	haveill   int
}

func MakeUser() UserData {
	return UserData{country: "",
		birth:     -1,
		gender:    -1,
		education: "",
		origin:    "",
		vaccine:   "",
		haveill:   0}
}

type UserSession struct {
	b              *Bot
	state          int
	userID, chatID int64
	userName       string
	userChan       chan string
	userData       UserData
}

func MakeSession(b *Bot, userID, chatID int64, userName string) *UserSession {
	return &UserSession{b, survey_default, userID, chatID, userName, make(chan string), MakeUser()}
}

func (s *UserSession) nextStep() {
	s.state++
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
	var res bool

	country := userData
	switch country {
	case Russia[1], Ukraine[1], Belarus[1], Kazakh[1]:
		s.userData.country = country
		res = true
	default:
		msg.Text = error_msg + "страна " + country + " не участвует в опросе"
		res = false
	}

	s.b.bot.Send(msg)
	return res
}

func (s *UserSession) askBirth_01() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_birth_msg
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getBirth_02(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var res bool

	year, _ := strconv.Atoi(userData)
	if 1920 < year && year < 2020 {
		s.userData.birth = year
		res = true
	} else {
		msg.Text = error_msg + userData + " не корректный год"
		res = false
	}

	s.b.bot.Send(msg)
	return res
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
	var res bool

	if userData == Male[1] || userData == Female[1] {
		s.userData.gender, _ = strconv.Atoi(userData)
		res = true
	} else {
		msg.Text = error_msg + "не корректный пол"
		res = false
	}

	s.b.bot.Send(msg)
	return res
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
	var res bool

	education := userData
	switch education {
	case School[1], College[1], University[1]:
		s.userData.education = education
		res = true
	default:
		msg.Text = error_msg + education + error_ans
		res = false
	}

	s.b.bot.Send(msg)
	return res
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
	var res bool

	origin := userData
	switch origin {
	case Nature[1], Human[1], Unknown[1]:
		s.userData.origin = origin
		res = true
	default:
		msg.Text = error_msg + origin + error_ans
		res = false
	}

	s.b.bot.Send(msg)
	return res
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
	var res bool

	vaccine := userData
	switch vaccine {
	case Helpful[1], Useless[1], Dangerous[1], Unknown[1]:
		s.userData.vaccine = vaccine
		res = true
	default:
		msg.Text = error_msg + vaccine + error_ans
		res = false
	}

	s.b.bot.Send(msg)
	return res
}

func (s *UserSession) askHaveIll_06() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	msg.Text = ask_haveill_msg
	msg.ReplyMarkup = haveillInlineKeyboard
	s.nextStep()

	s.b.bot.Send(msg)
}

func (s *UserSession) getHaveIll_07(userData string) bool {
	msg := tgbotapi.NewMessage(s.chatID, "")
	var res bool

	haveill := userData
	switch haveill {
	case Yes[1]:
		s.userData.haveill = 1
		res = true
	case No[1]:
		s.userData.haveill = 0
		res = true
	default:
		msg.Text = error_msg + haveill + error_ans
		res = false
	}

	s.b.bot.Send(msg)
	return res
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
			if s.getCountry_01(data) {
				s.askBirth_01()
			}
		case 2:
			if s.getBirth_02(data) {
				s.askGender_02()
			}
		case 3:
			if s.getGender_03(data) {
				s.askEducation_03()
			}
		case 4:
			if s.getEducation_04(data) {
				s.askOrigin_04()
			}
		case 5:
			if s.getOrigin_05(data) {
				s.askVaccine_05()
			}
		case 6:
			if s.getVaccine_06(data) {
				s.askHaveIll_06()
			}
		case 7:
			s.getHaveIll_07(data)
		default:
			s.writeResult(data)
			return
		}
	}
}
