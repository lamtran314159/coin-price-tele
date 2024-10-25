package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BinanceSpotWSURL   = "wss://stream.binance.com:9443/ws"
	BinanceFutureWSURL = "wss://fstream.binance.com/ws"
)

var conn *websocket.Conn

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
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@ticker"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
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
				log.Println("Error read message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
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
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
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
				log.Println("Error read message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
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
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
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
				log.Println("Error read message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
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
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
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
				log.Println("Error read message:", err)
				msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please provide a valid symbol")
				bot.Send(msg)
			}
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
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
