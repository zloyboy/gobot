package telegram

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const error_msg = "Произошла ошибка: "
const error_ans = " не корректный ответ"
const abort_msg = "Опрос завершен из-за ошибки. Попробуйте начать сначала, но не ранее чем через 10 секунд"

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

const ask_country_msg = "Укажите пожалуйста страну проживания:"

var countryInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readCountryFromDb() bool {
	country := b.dbase.ReadCaption("userCountry")
	if 4 <= len(country) {
		countryInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

const ask_birth_msg = "Введите пожалуйста год рождения"

const ask_gender_msg = "Укажите пожалуйста ваш пол"

var Male = [2]string{"Мужской", "1"}
var Female = [2]string{"Женский", "0"}
var genderInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Male[0], Male[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Female[0], Female[1]),
	),
)

const ask_education_msg = "Укажите пожалуйста ваше образование"

var educationInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readEducationFromDb() bool {
	education := b.dbase.ReadCaption("userEducation")
	if 3 <= len(education) {
		educationInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

const ask_vaccine_msg = "Считаете ли вы что существующие прививки (какие-то лучше, какие-то хуже) помогают предотвратить или облегчить болезнь?"

var vaccineInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readVaccineFromDb() bool {
	vaccine := b.dbase.ReadCaption("userVaccineOpinion")
	if 3 <= len(vaccine) {
		vaccineInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

const ask_origin_msg = "Считаете ли вы что новый коронавирус это естественный природный процесс или к его созданию причастны люди?"

var originInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readOriginFromDb() bool {
	origin := b.dbase.ReadCaption("userOriginOpinion")
	if 2 <= len(origin) {
		originInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

const ask_haveill_msg = "Считаете ли вы что переболели коронавирусом (возможно не один раз)?"
const ask_countill_msg = "Введите пожалуйста сколько раз вы переболели коронавирусом"
const ask_yearill_msg = "Введите год когда переболели "
const ask_monthill_msg = "Введите месяц когда переболели "
const ask_signill_msg = "По каким признакам вы определили тогда, что переболели коронавирусом?"

var signillInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readIllnessSignFromDb() bool {
	sign := b.dbase.ReadCaption("illnessSign")
	if 3 <= len(sign) {
		signillInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

const ask_degreeill_msg = "Насколько тяжело протекала болезнь?"

var degreeillInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readIllnessDegreeFromDb() bool {
	degree := b.dbase.ReadCaption("illnessDegree")
	if 6 <= len(degree) {
		degreeillInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
const ask_yearvac_msg = "Введите год когда вакцинировались "
const ask_monthvac_msg = "Введите месяц когда вакцинировались "
const ask_kindvac_msg = "Какую вакцину вводили?"

var kindvacInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readVaccineKindFromDb() bool {
	kind := b.dbase.ReadCaption("vaccineKind")
	if 4 <= len(kind) {
		kindvacInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

const ask_effectvac_msg = "Насколько сильными были побочные эффекты после вакцины?"

var effectvacInlineKeyboard tgbotapi.InlineKeyboardMarkup

func (b *Bot) readVaccineEffectFromDb() bool {
	effect := b.dbase.ReadCaption("vaccineEffect")
	if 3 <= len(effect) {
		effectvacInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
