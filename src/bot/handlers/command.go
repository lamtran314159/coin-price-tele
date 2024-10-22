package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"
	
	"encoding/json"
	"strconv"
	"time"
	"github.com/gorilla/websocket"

	"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
    binanceSpotWSURL   = "wss://stream.binance.com:9443/ws"
    binanceFutureWSURL = "wss://fstream.binance.com/ws"
)

var conn *websocket.Conn

func connectWebSocket(url string) error {
    var err error
    conn, _, err = websocket.DefaultDialer.Dial(url, nil)
    return err
}

func subscribeToStream(stream string) error {
    subscribeMsg := fmt.Sprintf(`{"method": "SUBSCRIBE", "params":["%s"], "id": 1}`, stream)
    return conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg))
}

func formatNumber(value interface{}) string {
    switch v := value.(type) {
    case string:
        if f, err := strconv.ParseFloat(v, 64); err == nil {
            return fmt.Sprintf("%g", f)
        }
    case float64:
        return fmt.Sprintf("%g", v)
    }
    return fmt.Sprintf("%v", value)
}

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
	//isPriceCommand := strings.HasPrefix(command, "/p ")
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
	case strings.HasPrefix(command, "/p_spot"):
		symbol := strings.TrimSpace(strings.TrimPrefix(command, "/p_spot"))
		go getSpotPrice(chatID, symbol, bot)
	case strings.HasPrefix(command, "/p_future"):
		symbol := strings.TrimSpace(strings.TrimPrefix(command, "/p_future"))
		go getFuturePrice(chatID, symbol, bot)
	case strings.HasPrefix(command, "/fundRate"):
		symbol := strings.TrimSpace(strings.TrimPrefix(command, "/fundRate"))
		go getFundingRate(chatID, symbol, bot)
	case strings.HasPrefix(command, "/fundRateCDown"):
		symbol := strings.TrimSpace(strings.TrimPrefix(command, "/fundRateCDown"))
		go getFundingRateCountdown(chatID, symbol, bot)
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

func getSpotPrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
    if err := connectWebSocket(binanceSpotWSURL); err != nil {
        log.Println("Lỗi kết nối WebSocket:", err)
        return
    }
    defer conn.Close()

    stream := strings.ToLower(symbol) + "@ticker"
    if err := subscribeToStream(stream); err != nil {
        log.Println("Lỗi đăng ký stream:", err)
        return
    }

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Lỗi đọc tin nhắn:", err)
            return
        }

        var data map[string]interface{}
        if err := json.Unmarshal(message, &data); err != nil {
            log.Println("Lỗi giải mã JSON:", err)
            continue
        }

        if price, ok := data["c"]; ok {
            formattedPrice := formatNumber(price)
            msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Giá Spot hiện tại của %s: %s", symbol, formattedPrice))
            bot.Send(msg)
            return
        }
    }
}

func getFuturePrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
    if err := connectWebSocket(binanceFutureWSURL); err != nil {
        log.Println("Lỗi kết nối WebSocket:", err)
        return
    }
    defer conn.Close()

    stream := strings.ToLower(symbol) + "@markPrice"
    if err := subscribeToStream(stream); err != nil {
        log.Println("Lỗi đăng ký stream:", err)
        return
    }

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Lỗi đọc tin nhắn:", err)
            return
        }

        var data map[string]interface{}
        if err := json.Unmarshal(message, &data); err != nil {
            log.Println("Lỗi giải mã JSON:", err)
            continue
        }

        if price, ok := data["p"]; ok {
            formattedPrice := formatNumber(price)
            msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Giá Future hiện tại của %s: %s", symbol, formattedPrice))
            bot.Send(msg)
            return
        }
    }
}

func getFundingRate(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
    if err := connectWebSocket(binanceFutureWSURL); err != nil {
        log.Println("Lỗi kết nối WebSocket:", err)
        return
    }
    defer conn.Close()

    stream := strings.ToLower(symbol) + "@markPrice"
    if err := subscribeToStream(stream); err != nil {
        log.Println("Lỗi đăng ký stream:", err)
        return
    }

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Lỗi đọc tin nhắn:", err)
            return
        }

        var data map[string]interface{}
        if err := json.Unmarshal(message, &data); err != nil {
            log.Println("Lỗi giải mã JSON:", err)
            continue
        }

        if rate, ok := data["r"]; ok {
            formattedRate := formatNumber(rate)
            msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding Rate hiện tại của %s: %s", symbol, formattedRate))
            bot.Send(msg)
            return
        }
    }
}

func getFundingRateCountdown(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
    if err := connectWebSocket(binanceFutureWSURL); err != nil {
        log.Println("Lỗi kết nối WebSocket:", err)
        return
    }
    defer conn.Close()

    stream := strings.ToLower(symbol) + "@markPrice"
    if err := subscribeToStream(stream); err != nil {
        log.Println("Lỗi đăng ký stream:", err)
        return
    }

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Lỗi đọc tin nhắn:", err)
            return
        }

        var data map[string]interface{}
        if err := json.Unmarshal(message, &data); err != nil {
            log.Println("Lỗi giải mã JSON:", err)
            continue
        }

        if nextFundingTime, ok := data["T"]; ok {
            nextFundingTimeInt := int64(nextFundingTime.(float64))
            countdown := time.Until(time.Unix(nextFundingTimeInt/1000, 0))
            msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Thời gian đến Funding Rate tiếp theo của %s: %v", symbol, countdown.Round(time.Second)))
            bot.Send(msg)
            return
        }
    }
}