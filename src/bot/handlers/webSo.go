package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/websocket"
)

const (
	telegramBotToken   = config.GetEnv("BOT_TOKEN")
	binanceSpotWSURL   = "wss://stream.binance.com:9443/ws"
	binanceFutureWSURL = "wss://fstream.binance.com/ws"
)

var (
	bot  *tgbotapi.BotAPI
	conn *websocket.Conn
)
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
func handleCommand(message *tgbotapi.Message) {
	command := message.Command()
	args := message.CommandArguments()

	switch command {
	case "start":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Chào mừng! Sử dụng các lệnh sau:\n/spot symbol\n/futureP symbol\n/funRate symbol\n/funRateCountDown symbol")
		bot.Send(msg)
	case "p_spot":
		go getSpotPrice(message.Chat.ID, args)
	case "p_future":
		go getFuturePrice(message.Chat.ID, args)
	case "fundRate":
		go getFundingRate(message.Chat.ID, args)
	case "fundRateCDown":
		go getFundingRateCountdown(message.Chat.ID, args)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Lệnh không hợp lệ. Sử dụng /start để xem danh sách lệnh.")
		bot.Send(msg)
	}
}

func connectWebSocket(url string) error {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	return err
}

func subscribeToStream(stream string) error {
	subscribeMsg := fmt.Sprintf(`{"method": "SUBSCRIBE", "params":["%s"], "id": 1}`, stream)
	return conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg))
}

func getSpotPrice(chatID int64, symbol string) {
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

func getFuturePrice(chatID int64, symbol string) {
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

func getFundingRate(chatID int64, symbol string) {
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

func getFundingRateCountdown(chatID int64, symbol string) {
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