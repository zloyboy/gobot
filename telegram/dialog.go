package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const error_msg = "Произошла ошибка: "
const error_ans = " не корректный ответ"

var Yes = [2]string{"Да", "Yes"}
var No = [2]string{"Нет", "No"}
var Unknown = [2]string{"Не знаю", "Unknown"}

const start_msg = "Этот опрос создан для независимого сбора информации по пандемии коронавируса в РФ и странах СНГ. " +
	"После прохождения опроса вам будет доступна собранная статистика."

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Start"),
	),
)

const ask_country_msg = "Укажите пожалуйста страну проживания:"

var Russia = [2]string{"Россия", "Russia"}
var Ukraine = [2]string{"Украина", "Ukraine"}
var Belarus = [2]string{"Беларусь", "Belarus"}
var Kazakh = [2]string{"Казахстан", "Kazakhstan"}
var countryInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Russia[0], Russia[1]),
		tgbotapi.NewInlineKeyboardButtonData(Ukraine[0], Ukraine[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Belarus[0], Belarus[1]),
		tgbotapi.NewInlineKeyboardButtonData(Kazakh[0], Kazakh[1]),
	),
)

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

var School = [2]string{"Среднее", "School"}
var College = [2]string{"Колледж", "College"}
var University = [2]string{"Университет", "University"}

var educationInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(School[0], School[1]),
		tgbotapi.NewInlineKeyboardButtonData(College[0], College[1]),
		tgbotapi.NewInlineKeyboardButtonData(University[0], University[1]),
	),
)

const ask_origin_msg = "Считаете ли вы что новый коронавирус это естественный природный процесс или к его созданию причастны люди?"

var Nature = [2]string{"Природа", "Nature"}
var Human = [2]string{"Люди", "Human"}

var originInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Nature[0], Nature[1]),
		tgbotapi.NewInlineKeyboardButtonData(Human[0], Human[1]),
		tgbotapi.NewInlineKeyboardButtonData(Unknown[0], Unknown[1]),
	),
)

const ask_vaccine_msg = "Считаете ли вы что существующие прививки (какие-то лучше, какие-то хуже) помогают предотвратить или облегчить болезнь?"

var Helpful = [2]string{"Помогают", "Helpful"}
var Useless = [2]string{"Бесполезны", "Useless"}
var Dangerous = [2]string{"Опасны", "Dangerous"}

var vaccineInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Helpful[0], Helpful[1]),
		tgbotapi.NewInlineKeyboardButtonData(Useless[0], Useless[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Dangerous[0], Dangerous[1]),
		tgbotapi.NewInlineKeyboardButtonData(Unknown[0], Unknown[1]),
	),
)

const ask_haveill_msg = "Считаете ли вы что переболели коронавирусом (возможно не один раз)?"

var haveillInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Yes[0], Yes[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(No[0], No[1]),
	),
)

var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
