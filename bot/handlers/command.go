package handlers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handle incoming messages (commands or regular text)
func HandleMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	user := message.From
	text := message.Text

	log.Printf("%s wrote: %s", user.FirstName, text)

	if strings.HasPrefix(text, "/") {
		handleCommand(message.Chat.ID, text, bot)
	} else if screaming {
		bot.Send(sendScreamedMessage(message))
	} else {
		bot.Send(copyMessage(message))
	}
}

// Handle commands (e.g., /scream, /whisper, /menu)
func handleCommand(chatId int64, command string, bot *tgbotapi.BotAPI) {
	switch command {
	case "/scream":
		screaming = true
		bot.Send(tgbotapi.NewMessage(chatId, "Screaming mode enabled."))
	case "/whisper":
		screaming = false
		bot.Send(tgbotapi.NewMessage(chatId, "Screaming mode disabled."))
	case "/menu":
		bot.Send(sendMenu(chatId))
	}
}

func sendMenu(chatId int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	return msg
}

func sendScreamedMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(message.Text))
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func copyMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	return msg
}
