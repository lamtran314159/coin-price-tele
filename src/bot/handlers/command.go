package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userTokens = struct {
	sync.RWMutex
	m map[int]string
}{m: make(map[int]string)}

// Handle incoming messages (commands or regular text)
func HandleMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	user := message.From
	text := message.Text

	log.Printf("%s wrote: %s", user.FirstName, text)

	if strings.HasPrefix(text, "/") {
		handleCommand(message.Chat.ID, text, bot, user)
	} else if screaming {
		_, err := bot.Send(sendScreamedMessage(message))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	} else {
		_, err := bot.Send(copyMessage(message))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

// Handle commands (e.g., /scream, /whisper, /menu)
func handleCommand(chatID int64, command string, bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	fmt.Println(user.ID)
	isPriceCommand := strings.HasPrefix(command, "/p ")
	switch {
	case command == "/help":
		_, err := bot.Send(tgbotapi.NewMessage(chatID, strings.Join(commandList, "\n")))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case command == "/start":
		token, err := services.AuthenticateUser(user.ID)
		if err != nil {
			_, err := bot.Send(tgbotapi.NewMessage(chatID, "Access denied."))
			if err != nil {
				log.Println("Error sending message:", err)
			}
			return
		}

		// Send a message showing access was granted
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Access granted. Your token is: %s", token))
		// Store the token
		userTokens.Lock()
		userTokens.m[int(user.ID)] = token
		userTokens.Unlock()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case command == "/scream":
		screaming = true
		_, err := bot.Send(tgbotapi.NewMessage(chatID, "Screaming mode enabled."))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case command == "/whisper":
		screaming = false
		_, err := bot.Send(tgbotapi.NewMessage(chatID, "Screaming mode disabled."))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case command == "/menu":
		_, err := bot.Send(sendMenu(chatID))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case command == "/protected":
		token := userTokens.m[int(user.ID)]
		response, err := services.ValidateToken(token)
		if err != nil {
			log.Println("Error validating token:", err)
			return
		}
		_, err = bot.Send(tgbotapi.NewMessage(chatID, response))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	// case isPriceCommand: // Handle the /p <symbol> command
	// 	symbol := strings.TrimSpace(command[3:])

	// 	if symbol == "" {
	// 		bot.Send(tgbotapi.NewMessage(chatID, "Please provide a symbol (e.g., /p eth)."))
	// 		return
	// 	}
	// 	//log.Printf("Symbol: %s", symbol)
	// 	price, exists := CryptoPrices[strings.ToUpper(symbol)+"USDT"]
	// 	//log.Printf(strings.ToUpper(symbol) + "USDT")
	// 	//log.Printf("Price: %f, Exists: %t", price, exists)
	// 	if !exists || price == 0 {
	// 		bot.Send(tgbotapi.NewMessage(chatID, "Price for "+symbol+" is not available yet. Please try again later."))
	// 		return
	// 	}
	// 	message := fmt.Sprintf("Current %s price: $%.4f", symbol, price)
	// 	_, err := bot.Send(tgbotapi.NewMessage(chatID, message))
	// 	if err != nil {
	// 		log.Println("Error sending message:", err)
	// 	}
	// }
	case command == "/p_spot":
		getSpotPrice(chatID, symbol)
	case command == "/p_future":
		getFuturePrice(chatID, symbol)
	case command == "/fundRate":
		getFundingRate(chatID, symbol)
	case command == "/fundRateCDown":
		getFundingRateCountdown(chatID, symbol)
	}
}

func sendMenu(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, firstMenu)
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
