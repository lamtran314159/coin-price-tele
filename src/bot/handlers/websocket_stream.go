package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BinanceSpotWSURL   = "wss://stream.binance.com:9443/ws"
	BinanceFutureWSURL = "wss://fstream.binance.com/ws"
)

const baseUrl = "https://hcmutssps.id.vn/api/vip1/get-kline"

var conn *websocket.Conn

// KlineData struct to hold data from API
type KlineData struct {
	Symbol             string `json:"symbol"`
	EventTime          string `json:"eventTime"`
	KlineStartTime     string `json:"klineStartTime"`
	KlineCloseTime     string `json:"klineCloseTime"`
	OpenPrice          string `json:"openPrice"`
	ClosePrice         string `json:"closePrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	NumberOfTrades     int    `json:"numberOfTrades"`
	BaseAssetVolume    string `json:"baseAssetVolume"`
	TakerBuyVolume     string `json:"takerBuyVolume"`
	TakerBuyBaseVolume string `json:"takerBuyBaseVolume"`
	Volume             string `json:"volume"`
}

func connectWebSocket(url string) error {
	var err error
	// Create a background context
	ctx := context.Background()
	conn, _, err = websocket.Dial(ctx, url, nil)
	return err
}

func subscribeToStream(stream string) error {
	// Create a background context for the write operation
	ctx := context.Background()
	subscribeMsg := fmt.Sprintf(`{"method": "SUBSCRIBE", "params":["%s"], "id": 1}`, stream)
	return conn.Write(ctx, websocket.MessageText, []byte(subscribeMsg))
}

func formatNumber(value interface{}, isFundingRate bool) string {
	switch v := value.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			// Format with 5 decimal places for funding rate, otherwise 4 decimal places for price
			if isFundingRate {
				return fmt.Sprintf("%.5f", f)
			}
			return fmt.Sprintf("%.4f", f)
		}
	case float64:
		// Format with 5 decimal places for funding rate, otherwise 4 decimal places for price
		if isFundingRate {
			return fmt.Sprintf("%.5f", v)
		}
		return fmt.Sprintf("%.4f", v)
	}
	// Default case if the value is neither string nor float64
	return fmt.Sprintf("%v", value)
}

func GetSpotPrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(BinanceSpotWSURL); err != nil {
		log.Println("Failed to connect to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	stream := strings.ToLower(symbol) + "@ticker"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Failed to subscribe to stream:", err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for {
		// Use timeoutCtx in the Read operation
		_, message, err := conn.Read(timeoutCtx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Println("WebSocket closed normally")
			} else {
				log.Println("Error reading message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol.")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Failed to unmarshal JSON:", err)
			continue
		}

		if price, ok := data["c"]; ok {
			formattedPrice := formatNumber(price, false)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Spot price | %s |  %s", symbol, formattedPrice))
			bot.Send(msg)
			return
		}
	}
}

func GetFuturePrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(BinanceFutureWSURL); err != nil {
		log.Println("Failed to connect to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Failed to subscribe to stream:", err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for {
		// Use timeoutCtx in the Read operation
		_, message, err := conn.Read(timeoutCtx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Println("WebSocket closed normally")
			} else {
				log.Println("Error reading message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Failed to unmarshal JSON:", err)
			continue
		}

		if price, ok := data["p"]; ok {
			formattedPrice := formatNumber(price, false)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Future price | %s |  %s", symbol, formattedPrice))
			bot.Send(msg)
			return
		}
	}
}

func GetFundingRate(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(BinanceFutureWSURL); err != nil {
		log.Println("Failed to connect to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Failed to subscribe to stream:", err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for {
		// Use timeoutCtx in the Read operation
		_, message, err := conn.Read(timeoutCtx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Println("WebSocket closed normally")
			} else {
				log.Println("Error reading message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Failed to unmarshal JSON:", err)
			continue
		}

		if fundingRate, ok := data["r"]; ok {
			formattedFundingRate := formatNumber(fundingRate, true)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding rate | %s |  %s", symbol, formattedFundingRate))
			bot.Send(msg)
			return
		}
	}
}

func GetFundingRateCountdown(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(BinanceFutureWSURL); err != nil {
		log.Println("Failed to connect to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Failed to subscribe to stream:", err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for {
		// Use timeoutCtx in the Read operation
		_, message, err := conn.Read(timeoutCtx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Println("WebSocket closed normally")
			} else {
				log.Println("Error reading message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Failed to unmarshal JSON:", err)
			continue
		}

		if nextFundingTime, ok := data["T"]; ok {
			nextFundingTimeMillis := int64(nextFundingTime.(float64))
			remainingTime := time.Until(time.UnixMilli(nextFundingTimeMillis)).Truncate(time.Second)

			// Format remaining time into hours, minutes, and seconds
			hours := int(remainingTime.Hours())
			minutes := int(remainingTime.Minutes()) % 60
			seconds := int(remainingTime.Seconds()) % 60

			formattedRemainingTime := fmt.Sprintf("%02dh%02dm%02ds", hours, minutes, seconds)

			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding rate countdown | %s |  %s", symbol, formattedRemainingTime))
			bot.Send(msg)
			return
		}
	}
}

// UserConnection stores request state for each user
type UserConnection struct {
	isStreaming bool
}

// Tạo map để lưu trạng thái người dùng
var userConnections = make(map[int64]*UserConnection)
var mapMutex = sync.Mutex{}

// fetchKlineData sends a GET request to the backend API with cookie for security
func fetchKlineData(symbol, interval, cookie string, chatID int64, bot *tgbotapi.BotAPI) {
	reqUrl := fmt.Sprintf("%s?symbols=%s&interval=%s", baseUrl, symbol, interval)
	log.Printf("API URL: %s", reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Request creation error: %v", err)))
		return
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch Kline data: %v", err)))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Received status code %d", resp.StatusCode)))
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		mapMutex.Lock()
		userConn := userConnections[chatID]
		if userConn == nil || !userConn.isStreaming {
			mapMutex.Unlock()
			return // Thoát vòng lặp khi isStreaming = false
		}
		mapMutex.Unlock()
		// Lấy dòng dữ liệu và loại bỏ tiền tố "data:"
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Giải mã JSON
		var klineData KlineData
		err := json.Unmarshal([]byte(line), &klineData)
		if err != nil {
			// bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error decoding response: %v", err)))
			continue
		}

		// Gửi dữ liệu đã giải mã đến người dùng Telegram
		klineMessage := fmt.Sprintf("Real-time Kline data:\nOpen: %s\nHigh: %s\nLow: %s\nClose: %s",
			klineData.OpenPrice, klineData.HighPrice, klineData.LowPrice, klineData.ClosePrice)
		bot.Send(tgbotapi.NewMessage(chatID, klineMessage))

		time.Sleep(time.Second) // Tránh spam tin nhắn
	}
}
