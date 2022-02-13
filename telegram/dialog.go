package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	error_msg = "Произошла ошибка: "
	error_ans = " не корректный ответ"
	abort_msg = "Опрос завершен из-за ошибки. Попробуйте начать сначала, но не ранее чем через 10 секунд"
	start_msg = "Этот опрос создан для независимого сбора информации по пандемии коронавируса в РФ и странах СНГ. " +
		"После прохождения опроса вам будет доступна собранная статистика."
	repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
)

type BaseQuestion struct {
	ask string
	key interface{}
	min int
	max int
}

var (
	baseQuestion = []BaseQuestion{
		{"Укажите пожалуйста страну проживания", nil, 0, 3},
		{"Введите пожалуйста год рождения", nil, 1920, 2020},
		{"Укажите пожалуйста ваш пол", nil, 0, 1},
		{"Укажите пожалуйста ваше образование", nil, 0, 2},
		{"Считаете ли вы что существующие прививки (российские и ииностранные) какие-то лучше, какие-то хуже, но помогают предотвратить или облегчить болезнь?", nil, -1, 2},
		{"Считаете ли вы что новый коронавирус это естественный природный процесс или к его созданию причастны люди?", nil, -1, 1},
		{"Считаете ли вы что переболели коронавирусом (возможно не один раз)?", nil, 0, 0},
	}

	Yes     = [2]string{"Да", "1"}
	No      = [2]string{"Нет", "0"}
	Unknown = [2]string{"Не знаю", "-1"}

	yesnoInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Yes[0], Yes[1]),
			tgbotapi.NewInlineKeyboardButtonData(No[0], No[1]),
		),
	)

	numInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
	)
)

func (b *Bot) readYearFromDb() bool {
	year := b.dbase.ReadCaption("year", 2020) // years start with 2020
	if 3 <= len(year) {
		yearInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(year[0][0], year[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(year[1][0], year[1][1]),
				tgbotapi.NewInlineKeyboardButtonData(year[2][0], year[2][1]),
			),
		)
		illQuestion[sst_year].key = yearInlineKeyboard
		vacQuestion[sst_year].key = yearInlineKeyboard
		return true
	}
	log.Print("Couldn't read years")
	return false
}

func (b *Bot) readMonthFromDb() bool {
	month := b.dbase.ReadCaption("month", 1) // months start with 01
	if len(month) == 12 {
		monthInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(month[0][0], month[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[1][0], month[1][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[2][0], month[2][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[3][0], month[3][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(month[4][0], month[4][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[5][0], month[5][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[6][0], month[6][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[7][0], month[7][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(month[8][0], month[8][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[9][0], month[9][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[10][0], month[10][1]),
				tgbotapi.NewInlineKeyboardButtonData(month[11][0], month[11][1]),
			),
		)
		illQuestion[sst_month].key = monthInlineKeyboard
		vacQuestion[sst_month].key = monthInlineKeyboard
		return true
	}
	log.Print("Couldn't read months")
	return false
}

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

func (b *Bot) readVaccOpinionFromDb() bool {
	vaccine := b.dbase.ReadCaption("userVaccineOpinion")
	if 3 <= len(vaccine) {
		baseQuestion[st_vacc_opin].key = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(vaccine[0][0], vaccine[0][1]),
				tgbotapi.NewInlineKeyboardButtonData(vaccine[1][0], vaccine[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(vaccine[2][0], vaccine[2][1]),
				tgbotapi.NewInlineKeyboardButtonData(Unknown[0], Unknown[1]),
			),
		)
		return true
	}
	log.Print("Couldn't read vaccine")
	return false
}

func (b *Bot) readOrgnOpinionFromDb() bool {
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

const ask_countill_msg = "Введите пожалуйста сколько раз вы переболели коронавирусом"

type SubQuestion struct {
	ask  string
	key  interface{}
	min  int
	max  int
	many bool
}

var illQuestion = []SubQuestion{
	{"Введите год когда переболели ", nil, 2020, 2022, true},
	{"Введите месяц когда переболели ", nil, 1, 12, true},
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
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[1][0], degree[1][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[2][0], degree[2][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[3][0], degree[3][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(degree[4][0], degree[4][1]),
			),
			tgbotapi.NewInlineKeyboardRow(
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
	{"Введите год когда вакцинировались ", nil, 2020, 2022, true},
	{"Введите месяц когда вакцинировались ", nil, 1, 12, true},
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
