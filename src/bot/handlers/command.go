package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		parts := strings.Fields(text)
		command := parts[0]
		args := parts[1:]
		handleCommand(message.Chat.ID, command, args, bot, user)
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
func handleCommand(chatID int64, command string, args []string, bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	fmt.Println(user.ID)
	switch command {
	case "/help":
		_, err := bot.Send(tgbotapi.NewMessage(chatID, strings.Join(commandList, "\n")))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/start":
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
	case "/scream":
		screaming = true
		_, err := bot.Send(tgbotapi.NewMessage(chatID, "Screaming mode enabled."))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/whisper":
		screaming = false
		_, err := bot.Send(tgbotapi.NewMessage(chatID, "Screaming mode disabled."))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/kline":
		fmt.Println(args)

		if len(args) < 2 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /kline <symbol> <interval> [limit]")
			bot.Send(msg)
			return
		}

		symbol := args[0]
		interval := args[1]

		// Set default limit if not provided
		limit := 5 // Default value
		if len(args) == 3 {
			parsedLimit, err := strconv.Atoi(args[2])
			if err == nil {
				limit = parsedLimit
			}
		}

		data, err := getKlineData(symbol, interval, limit) // Pass parameters as needed
		if err != nil {
			// If there's an error, send the error message back to the user
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error fetching Kline data: "+err.Error()))
		} else {
			// If there's no error, send the Kline data
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, data))
		}
	case "/menu":
		_, err := bot.Send(sendMenu(chatID))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/protected":
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
	}
	// case "/spot-lower":

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

func getKlineData(symbol string, interval string, limit int) (string, error) {
	// Define the API endpoint
	apiURL := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)

	// Make the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Return the raw JSON response as a string
	return string(body), nil
}
