package telegram

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const error_msg = "Произошла ошибка: "
const error_ans = " не корректный ответ"
const abort_msg = "Опрос завершен из-за ошибки. Попробуйте начать сначала, но не ранее чем через 10 секунд"

type BaseQuestion struct {
	ask string
	key interface{}
	min int
	max int
}

var baseQuestion = []BaseQuestion{
	{"Укажите пожалуйста страну проживания", nil, 0, 3},
	{"Введите пожалуйста год рождения", nil, 1920, 2020},
	{"Укажите пожалуйста ваш пол", nil, 0, 1},
	{"Укажите пожалуйста ваше образование", nil, 0, 2},
	{"Считаете ли вы что существующие прививки (какие-то лучше, какие-то хуже) помогают предотвратить или облегчить болезнь?", nil, -1, 2},
	{"Считаете ли вы что новый коронавирус это естественный природный процесс или к его созданию причастны люди?", nil, -1, 1},
	{"Считаете ли вы что переболели коронавирусом (возможно не один раз)?", nil, 0, 0},
}

func nTimes(time, count int) string {
	if 1 < count {
		return strconv.Itoa(time) + "й раз"
	} else {
		return ""
	}
}

var Yes = [2]string{"Да", "Yes"}
var No = [2]string{"Нет", "No"}
var Unknown = [2]string{"Не знаю", "Unknown"}

var yesnoInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Yes[0], Yes[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(No[0], No[1]),
	),
)

const year2020 = "2020"
const year2021 = "2021"
const year2022 = "2022"

var yearInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(year2020, year2020),
		tgbotapi.NewInlineKeyboardButtonData(year2021, year2021),
		tgbotapi.NewInlineKeyboardButtonData(year2022, year2022),
	),
)

var January = [2]string{"Январь", "1"}
var February = [2]string{"Февраль", "2"}
var March = [2]string{"Март", "3"}
var April = [2]string{"Апрель", "4"}
var May = [2]string{"Май", "5"}
var June = [2]string{"Июнь", "6"}
var July = [2]string{"Июль", "7"}
var August = [2]string{"Август", "8"}
var September = [2]string{"Сентябрь", "9"}
var October = [2]string{"Октябрь", "10"}
var November = [2]string{"Ноябрь", "11"}
var December = [2]string{"Декабрь", "12"}

var monthInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(January[0], January[1]),
		tgbotapi.NewInlineKeyboardButtonData(February[0], February[1]),
		tgbotapi.NewInlineKeyboardButtonData(March[0], March[1]),
		tgbotapi.NewInlineKeyboardButtonData(April[0], April[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(May[0], May[1]),
		tgbotapi.NewInlineKeyboardButtonData(June[0], June[1]),
		tgbotapi.NewInlineKeyboardButtonData(July[0], July[1]),
		tgbotapi.NewInlineKeyboardButtonData(August[0], August[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(September[0], September[1]),
		tgbotapi.NewInlineKeyboardButtonData(October[0], October[1]),
		tgbotapi.NewInlineKeyboardButtonData(November[0], November[1]),
		tgbotapi.NewInlineKeyboardButtonData(December[0], December[1]),
	),
)

const start_msg = "Этот опрос создан для независимого сбора информации по пандемии коронавируса в РФ и странах СНГ. " +
	"После прохождения опроса вам будет доступна собранная статистика."

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Start"),
	),
)

func (b *Bot) readCountryFromDb() bool {
	country := b.dbase.ReadCaption("userCountry")
	if 4 <= len(country) {
		baseQuestion[st_country].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(country[0][0], country[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(country[1][0], country[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(country[2][0], country[2][1]),
				tgbotapi.NewInlineKeyboardButtonData(country[3][0], country[3][1]),
			),
		)
		return true
	}
	log.Print("Couldn't read countries")
	return false
}

var Male = [2]string{"Мужской", "1"}
var Female = [2]string{"Женский", "0"}

func (b *Bot) setKeyboards() {
	baseQuestion[st_gender].key = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Male[0], Male[1]),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Female[0], Female[1]),
		),
	)

	baseQuestion[st_have_ill].key = yesnoInlineKeyboard
}

func (b *Bot) readEducationFromDb() bool {
	education := b.dbase.ReadCaption("userEducation")
	if 3 <= len(education) {
		baseQuestion[st_education].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(education[0][0], education[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(education[1][0], education[1][1]),
				tgbotapi.NewInlineKeyboardButtonData(education[2][0], education[2][1]),
			),
		)
		return true
	}
	log.Print("Couldn't read education")
	return false
}

func (b *Bot) readVaccineFromDb() bool {
	vaccine := b.dbase.ReadCaption("userVaccineOpinion")
	if 3 <= len(vaccine) {
		baseQuestion[st_vacc_opin].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(vaccine[0][0], vaccine[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(vaccine[1][0], vaccine[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(vaccine[2][0], vaccine[2][1]),
				tgbotapi.NewInlineKeyboardButtonData(Unknown[0], "-1"),
			),
		)
		return true
	}
	log.Print("Couldn't read vaccine")
	return false
}

func (b *Bot) readOriginFromDb() bool {
	origin := b.dbase.ReadCaption("userOriginOpinion")
	if 2 <= len(origin) {
		baseQuestion[st_orgn_opin].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(origin[0][0], origin[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(origin[1][0], origin[1][1]),
				tgbotapi.NewInlineKeyboardButtonData(Unknown[0], Unknown[1]),
			),
		)
		return true
	}
	log.Print("Couldn't read origin")
	return false
}

type SubQuestion struct {
	ask  string
	key  interface{}
	min  int
	max  int
	many bool
}

const ask_countill_msg = "Введите пожалуйста сколько раз вы переболели коронавирусом"

var illQuestion = []SubQuestion{
	{"Введите год когда переболели ", yearInlineKeyboard, 2020, 2022, false},
	{"Введите месяц когда переболели ", monthInlineKeyboard, 1, 12, false},
	{"По каким признакам вы определили тогда, что переболели коронавирусом?", nil, 0, 2, false},
	{"Насколько тяжело протекала болезнь?", nil, 0, 5, false},
}

func (b *Bot) readIllnessSignFromDb() bool {
	sign := b.dbase.ReadCaption("illnessSign")
	if 3 <= len(sign) {
		illQuestion[sst_sign].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(sign[0][0], sign[0][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(sign[1][0], sign[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(sign[2][0], sign[2][1]),
			),
		)
		return true
	}
	log.Print("Couldn't read illness signs")
	return false
}

func (b *Bot) readIllnessDegreeFromDb() bool {
	degree := b.dbase.ReadCaption("illnessDegree")
	if 6 <= len(degree) {
		illQuestion[sst_degree].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[0][0], degree[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(degree[1][0], degree[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[2][0], degree[2][1]),
				tgbotapi.NewInlineKeyboardButtonData(degree[3][0], degree[3][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[4][0], degree[4][1]),
				tgbotapi.NewInlineKeyboardButtonData(degree[5][0], degree[5][1]),
			),
		)
		return true
	}
	log.Print("Couldn't read illness degrees")
	return false
}

const ask_havevac_msg = "Вы делали вакцинацию от коронавируса?"
const ask_countvac_msg = "Сколько раз вы вакцинировались?\n(Два укола Спутник-V считаются одним разом)"

var vacQuestion = []SubQuestion{
	{"Введите год когда вакцинировались ", yearInlineKeyboard, 2020, 2022, false},
	{"Введите месяц когда вакцинировались ", monthInlineKeyboard, 1, 12, false},
	{"Какую вакцину вводили?", nil, 0, 3, false},
	{"Насколько сильными были побочные эффекты после вакцины?", nil, 0, 2, false},
}

func (b *Bot) readVaccineKindFromDb() bool {
	kind := b.dbase.ReadCaption("vaccineKind")
	if 4 <= len(kind) {
		vacQuestion[sst_kind].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(kind[0][0], kind[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(kind[1][0], kind[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(kind[2][0], kind[2][1]),
				tgbotapi.NewInlineKeyboardButtonData(kind[3][0], kind[3][1]),
			),
		)
		return true
	}
	log.Print("Couldn't read vaccine kind")
	return false
}

func (b *Bot) readVaccineEffectFromDb() bool {
	effect := b.dbase.ReadCaption("vaccineEffect")
	if 3 <= len(effect) {
		vacQuestion[sst_effect].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(effect[0][0], effect[0][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(effect[1][0], effect[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(effect[2][0], effect[2][1]),
			),
		)
		return true
	}
	log.Print("Couldn't read vaccine effect")
	return false
}

var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
