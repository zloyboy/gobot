package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/zloyboy/gobot/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	st_country = user.Idx_country
	st_birth   = iota
	st_gender
	st_education
	st_vacc_opin
	st_orgn_opin
	st_have_ill
	st_get_have_ill
	st_get_count_ill
	st_illness
	st_get_have_vac
	st_get_count_vac
	st_vaccination
	st_check_result
)

const (
	sst_year   = user.Idx_year
	sst_month  = user.Idx_month
	sst_sign   = user.Idx_sign
	sst_degree = user.Idx_degree
	sst_kind   = user.Idx_kind
	sst_effect = user.Idx_effect
	sst_next   = 4
)

type UserSession struct {
	b              *Bot
	state          int
	subState       int
	count          int
	userID, chatID int64
	userData       user.UserData
	userSub        [4]int
}

func makeUserSession(b *Bot, userID, chatID int64) *UserSession {
	return &UserSession{b,
		0, 0, 0,
		userID, chatID,
		user.MakeUser(),
		user.MakeSubUser()}
}

func (s *UserSession) startSurvey() bool {
	if s.b.dbase.ExistId(s.userID) {
		s.sendStatic("Вы уже приняли участие в опросе")
		return false
	} else {
		// new user - start survey
		msg := tgbotapi.NewMessage(s.chatID, start_msg)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		s.b.bot.Send(msg)

		// first question
		msg = tgbotapi.NewMessage(s.chatID, baseQuestion[st_country].ask)
		msg.ReplyMarkup = baseQuestion[st_country].key
		s.b.bot.Send(msg)

		return true
	}
}

func (s *UserSession) sendStatic(header string) {
	msg := tgbotapi.NewMessage(s.chatID, header)
	s.b.bot.Send(msg)

	all := tgbotapi.NewPhoto(s.chatID, s.b.stat.MakeCommonChart())
	s.b.bot.Send(all)
	ill := tgbotapi.NewPhoto(s.chatID, s.b.stat.MakeChartIll())
	s.b.bot.Send(ill)
	vac := tgbotapi.NewPhoto(s.chatID, s.b.stat.MakeChartVac())
	s.b.bot.Send(vac)

	msg = tgbotapi.NewMessage(s.chatID, repeat_msg)
	msg.ReplyMarkup = startKeyboard
	s.b.bot.Send(msg)
}

func (s *UserSession) nextStep() {
	s.state++
}

func (s *UserSession) resetSubStep() {
	s.subState = 0
}

func (s *UserSession) nextSubStep() {
	s.subState++
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
	case st_birth:
		if s.count == 0 {
			s.sendQuestion(baseQuestion[st_birth].ask, baseQuestion[st_birth].key)
			s.state = st_birth // fix previous increment
			s.count = 1000
			s.userData.Base[st_birth] = 0
		}
	case st_gender, st_education, st_vacc_opin, st_orgn_opin, st_have_ill:
		s.sendQuestion(baseQuestion[s.state].ask, baseQuestion[s.state].key)
	case st_get_have_ill:
		if 0 < s.userData.CountIll {
			s.sendQuestion(ask_countill_msg, countInlineKeyboard)
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
			s.sendQuestion(ask_countvac_msg, countInlineKeyboard)
		} else {
			s.nextStep()
			s.nextStep()
			s.checkResult()
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
				s.checkResult()
			}
		}
	case st_check_result:
		return false
	}

	return true
}

func (s *UserSession) getAnswer(userData string) bool {
	var ok = true

	switch s.state {
	case st_country:
		val, err := strconv.Atoi(userData)
		if err != nil || val < baseQuestion[st_country].min || baseQuestion[st_country].max < val {
			ok = false
		} else {
			s.userData.Base[st_country] = val
			s.state = st_birth
		}
	case st_birth:
		digit, err := strconv.Atoi(userData)
		if err != nil {
			ok = false
		} else if 1920 <= digit && digit <= 2020 {
			s.userData.Base[st_birth] = digit
			s.count = 0
			s.state = st_gender
		} else if digit < 0 || 9 < digit || s.count == 1000 && digit != 1 && digit != 2 {
			ok = false
		} else {
			s.userData.Base[st_birth] += digit * s.count
			s.count /= 10
			if s.count == 0 {
				if 1920 <= s.userData.Base[st_birth] && s.userData.Base[st_birth] <= 2020 {
					s.b.bot.Send(tgbotapi.NewMessage(s.chatID, strconv.Itoa(s.userData.Base[st_birth])))
					s.state = st_gender
				} else {
					ok = false
				}
			}
		}
	case st_gender, st_education, st_vacc_opin, st_orgn_opin, st_have_ill:
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
	case st_check_result:
		val, err := strconv.Atoi(userData)
		if err != nil || val != 1 {
			ok = false
		}
	}

	return ok
}

func (s *UserSession) checkResult() {
	msg := tgbotapi.NewMessage(s.chatID, "")

	idx_country := s.userData.Base[st_country]
	birth_year := s.userData.Base[st_birth]
	idx_gender := s.userData.Base[st_gender]
	idx_education := s.userData.Base[st_education]
	idx_vacc_opin := s.userData.Base[st_vacc_opin]
	idx_orgn_opin := s.userData.Base[st_orgn_opin]

	msg.Text = "\nСтрана: " + country[idx_country][0] +
		"\nГод рождения: " + strconv.Itoa(birth_year) +
		"\nПол: " + gender[idx_gender][0] +
		"\nОбразование: " + education[idx_education][0] +
		"\nПрививки: " + vaccine[idx_vacc_opin][0] +
		"\nПроисхождение вируса: " + origin[idx_orgn_opin][0] +
		"\n--------------------"

	if 0 < s.userData.CountIll {
		msg.Text += fmt.Sprintf("\nБолел(а) %d раз(а)", s.userData.CountIll)
		for i := 0; i < s.userData.CountIll; i++ {
			ill_year := s.userData.Ill[i][sst_year]
			idx_month := s.userData.Ill[i][sst_month] - 1
			msg.Text += fmt.Sprintf("\n%dй раз: %d %s", i+1, ill_year, month[idx_month][0])
			idx_sign := s.userData.Ill[i][sst_sign]
			msg.Text += "\nПризнаки: " + ill_sign[idx_sign][0]
			idx_degree := s.userData.Ill[i][sst_degree]
			msg.Text += "\nТяжесть: " + ill_degree[idx_degree][0]
		}
	} else {
		msg.Text += "\nНе болел(а)"
	}
	msg.Text += "\n--------------------"

	if 0 < s.userData.CountVac {
		msg.Text += fmt.Sprintf("\nВакцинирован(а) %d раз(а)", s.userData.CountVac)
		for i := 0; i < s.userData.CountVac; i++ {
			vac_year := s.userData.Vac[i][sst_year]
			idx_month := s.userData.Vac[i][sst_month] - 1
			msg.Text += fmt.Sprintf("\n%dй раз: %d %s", i+1, vac_year, month[idx_month][0])
			idx_kind := s.userData.Vac[i][sst_kind]
			msg.Text += "\nВакцина: " + vac_kind[idx_kind][0]
			idx_effect := s.userData.Vac[i][sst_effect]
			msg.Text += "\nПобочки: " + vac_effect[idx_effect][0]
		}
	} else {
		msg.Text += "\nНе вакцинирован(а)"
	}
	msg.Text += "\n--------------------" +
		"\nВсе правильно?"

	msg.ReplyMarkup = cancelInlineKeyboard
	s.b.bot.Send(msg)
	s.nextStep()
}

func (s *UserSession) writeResult() {
	idx_country := s.userData.Base[st_country]
	birth_year := s.userData.Base[st_birth]
	idx_gender := s.userData.Base[st_gender]
	idx_education := s.userData.Base[st_education]
	idx_vacc_opin := s.userData.Base[st_vacc_opin]
	idx_orgn_opin := s.userData.Base[st_orgn_opin]
	log.Printf("insert id %d, country %s, birth %d, gender %d, edu %s, vaccine %s, origin %s, countIll %d, countVac %d",
		s.userID, country[idx_country][0], birth_year, idx_gender, education[idx_education][0],
		vaccine[idx_vacc_opin][0], origin[idx_orgn_opin][0], s.userData.CountIll, s.userData.CountVac)

	s.b.dbase.Insert(s.userID,
		time.Now().Local().Format("2006-01-02 15:04:05"),
		s.userData)
	s.b.stat.RefreshStatic(s.userData)
	s.sendStatic(thank_msg)
}

func (s *UserSession) abort(reason string) {
	msg := tgbotapi.NewMessage(s.chatID, reason+again_msg)
	msg.ReplyMarkup = startKeyboard
	s.b.bot.Send(msg)
}

const userTout = 10

func RunSurvey(b *Bot, userID, chatID int64, uchan *Channel) {
	tout := time.NewTimer(userTout * time.Second)

	s := makeUserSession(b, userID, chatID)

	defer func() {
		s.b.tout.SetTimer(s.userID)
		log.Printf("Exit user %d", s.userID)
		delete(s.b.uchan, s.userID)
		tout.Stop()
	}()

	if s.startSurvey() {
		for {
			select {
			case <-uchan.stop:
				s.abort(stop_msg)
				return
			case <-tout.C:
				s.abort(tout_msg)
				return
			case data := <-uchan.data:
				if !s.getAnswer(data) {
					s.abort(error_msg)
					return
				}
				if !tout.Stop() {
					<-tout.C
				}
				tout.Reset(userTout * time.Second)
				if !s.sendRequest() {
					s.writeResult()
				}
			}
		}
	}
}
