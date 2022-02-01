package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const error_msg = "Произошла ошибка: "
const error_ans = " не корректный ответ"
const abort_msg = "Опрос завершен из-за ошибки. Попробуйте начать сначала, но не ранее чем через 10 секунд"

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
const ask_countill_msg = "Введите пожалуйста сколько раз вы переболели коронавирусом"
const ask_yearill_msg = "Введите год когда переболели "
const ask_monthill_msg = "Введите месяц когда переболели "
const ask_signill_msg = "По каким признакам вы определили тогда, что переболели коронавирусом?"

var Medic = [2]string{"Есть медицинская справка", "medic"}
var Test = [2]string{"Есть тест с наличием антител", "test"}
var Symptom = [2]string{"По характерным симптомам", "symptom"}

var signillInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Medic[0], Medic[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Test[0], Test[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(Symptom[0], Symptom[1]),
	),
)

const ask_degreeill_msg = "Насколько тяжело протекала болезнь?"

var HospIvl = [2]string{"Лежал(а) под ИВЛ", "ivl"}
var Hospital = [2]string{"Лежал(а) в больнице", "hosp"}
var HomeHard = [2]string{"Лежал(а) дома, тяжело", "hard"}
var HomeEasy = [2]string{"Лежал(а) дома, легко", "easy"}
var OnFoot = [2]string{"Перенес(ла) на ногах", "onfoot"}
var NoSymptom = [2]string{"Перенес(ла) без симптомов", "nosymptom"}

var degreeillInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(HospIvl[0], HospIvl[1]),
		tgbotapi.NewInlineKeyboardButtonData(Hospital[0], Hospital[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(HomeHard[0], HomeHard[1]),
		tgbotapi.NewInlineKeyboardButtonData(HomeEasy[0], HomeEasy[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(OnFoot[0], OnFoot[1]),
		tgbotapi.NewInlineKeyboardButtonData(NoSymptom[0], NoSymptom[1]),
	),
)

const ask_havevac_msg = "Вы делали вакцинацию от коронавируса?"
const ask_countvac_msg = "Сколько раз вы вакцинировались?\n(Два укола Спутник-V считаются одним разом)"
const ask_yearvac_msg = "Введите год когда вакцинировались "
const ask_monthvac_msg = "Введите месяц когда вакцинировались "
const ask_kindvac_msg = "Какую вакцину вводили?"

var SputnikV = [2]string{"Спутник-V (два укола)", "sputnik-v"}
var SputnikL = [2]string{"Спутник-Лайт (один укол)", "sputnik-l"}
var EpiVac = [2]string{"ЭпиВакКорона", "epivac"}
var Kovivak = [2]string{"КовиВак", "sputnik-v"}

var kindvacInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(SputnikV[0], SputnikV[1]),
		tgbotapi.NewInlineKeyboardButtonData(SputnikL[0], SputnikL[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(EpiVac[0], EpiVac[1]),
		tgbotapi.NewInlineKeyboardButtonData(Kovivak[0], Kovivak[1]),
	),
)

const ask_effectvac_msg = "Насколько сильными были побочные эффекты после вакцины?"

var HardEffect = [2]string{"Сильные: температура, головная боль и т.п.", "hard"}
var MediumEffect = [2]string{"Средние: боль в руке, аллергия и т.п.", "medium"}
var EasyEffect = [2]string{"Слабые или никаких проявлений", "easy"}

var effectvacInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(HardEffect[0], HardEffect[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(MediumEffect[0], MediumEffect[1]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(EasyEffect[0], EasyEffect[1]),
	),
)

var repeat_msg = "\nДля повторного показа статистики введите любой текст или нажмите Start, но не ранее чем через 10 секунд"
