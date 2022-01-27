package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var numericInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ages[0], ages[0]),
		tgbotapi.NewInlineKeyboardButtonData(ages[1], ages[1]),
		tgbotapi.NewInlineKeyboardButtonData(ages[2], ages[2]),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ages[3], ages[3]),
		tgbotapi.NewInlineKeyboardButtonData(ages[4], ages[4]),
		tgbotapi.NewInlineKeyboardButtonData(ages[5], ages[5]),
	),
)

var numericResKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "1"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "0"),
	),
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Start"),
	),
)
