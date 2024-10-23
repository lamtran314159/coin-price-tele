package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Menu texts
const (
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."
)

// Button texts
const (
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"
)

var (
	screaming = false
)

var (
	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

var commandList = []string{
	"/start - Authenticate and start using the bot",
	"/scream - Enable screaming mode",
	"/whisper - Disable screaming mode",
	"/menu - Show menu with buttons",
	"/help - Show available commands",
	"/kline <symbol> <interval> [limit] [startTime] [endTime] - Get Kline data for a symbol",
}
