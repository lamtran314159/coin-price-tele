package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"strings"
// 	"time"
// 	"strconv"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"github.com/gorilla/websocket"
// )

// const (
// 	telegramBotToken   = config.GetEnv("BOT_TOKEN")
// 	binanceSpotWSURL   = "wss://stream.binance.com:9443/ws"
// 	binanceFutureWSURL = "wss://fstream.binance.com/ws"
// )

// func getSpotPrice(chatID int64, symbol string) {
// 	if err := connectWebSocket(binanceSpotWSURL); err != nil {
// 		log.Println("Lỗi kết nối WebSocket:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	stream := strings.ToLower(symbol) + "@ticker"
// 	if err := subscribeToStream(stream); err != nil {
// 		log.Println("Lỗi đăng ký stream:", err)
// 		return
// 	}

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Lỗi đọc tin nhắn:", err)
// 			return
// 		}

// 		var data map[string]interface{}
// 		if err := json.Unmarshal(message, &data); err != nil {
// 			log.Println("Lỗi giải mã JSON:", err)
// 			continue
// 		}

// 		if price, ok := data["c"]; ok {
// 			formattedPrice := formatNumber(price)
// 			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Giá Spot hiện tại của %s: %s", symbol, formattedPrice))
// 			bot.Send(msg)
// 			return
// 		}
// 	}
// }

// func getFuturePrice(chatID int64, symbol string) {
// 	if err := connectWebSocket(binanceFutureWSURL); err != nil {
// 		log.Println("Lỗi kết nối WebSocket:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	stream := strings.ToLower(symbol) + "@markPrice"
// 	if err := subscribeToStream(stream); err != nil {
// 		log.Println("Lỗi đăng ký stream:", err)
// 		return
// 	}

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Lỗi đọc tin nhắn:", err)
// 			return
// 		}

// 		var data map[string]interface{}
// 		if err := json.Unmarshal(message, &data); err != nil {
// 			log.Println("Lỗi giải mã JSON:", err)
// 			continue
// 		}

// 		if price, ok := data["p"]; ok {
// 			formattedPrice := formatNumber(price)
// 			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Giá Future hiện tại của %s: %s", symbol, formattedPrice))
// 			bot.Send(msg)
// 			return
// 		}
// 	}
// }

// func getFundingRate(chatID int64, symbol string) {
// 	if err := connectWebSocket(binanceFutureWSURL); err != nil {
// 		log.Println("Lỗi kết nối WebSocket:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	stream := strings.ToLower(symbol) + "@markPrice"
// 	if err := subscribeToStream(stream); err != nil {
// 		log.Println("Lỗi đăng ký stream:", err)
// 		return
// 	}

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Lỗi đọc tin nhắn:", err)
// 			return
// 		}

// 		var data map[string]interface{}
// 		if err := json.Unmarshal(message, &data); err != nil {
// 			log.Println("Lỗi giải mã JSON:", err)
// 			continue
// 		}

// 		if rate, ok := data["r"]; ok {
// 			formattedRate := formatNumber(rate)
// 			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding Rate hiện tại của %s: %s", symbol, formattedRate))
// 			bot.Send(msg)
// 			return
// 		}
// 	}
// }

// func getFundingRateCountdown(chatID int64, symbol string) {
// 	if err := connectWebSocket(binanceFutureWSURL); err != nil {
// 		log.Println("Lỗi kết nối WebSocket:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	stream := strings.ToLower(symbol) + "@markPrice"
// 	if err := subscribeToStream(stream); err != nil {
// 		log.Println("Lỗi đăng ký stream:", err)
// 		return
// 	}

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Lỗi đọc tin nhắn:", err)
// 			return
// 		}

// 		var data map[string]interface{}
// 		if err := json.Unmarshal(message, &data); err != nil {
// 			log.Println("Lỗi giải mã JSON:", err)
// 			continue
// 		}

// 		if nextFundingTime, ok := data["T"]; ok {
// 			nextFundingTimeInt := int64(nextFundingTime.(float64))
// 			countdown := time.Until(time.Unix(nextFundingTimeInt/1000, 0))
// 			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Thời gian đến Funding Rate tiếp theo của %s: %v", symbol, countdown.Round(time.Second)))
// 			bot.Send(msg)
// 			return
// 		}
// 	}
// }