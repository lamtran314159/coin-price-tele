package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coder/websocket"


	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


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
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceSpotWSURL); err != nil {
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@ticker"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Error read message:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
			continue
		}

		if price, ok := data["c"]; ok {
			formattedPrice := formatNumber(price)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Spot price | %s |  %s", symbol, formattedPrice))
			bot.Send(msg)
			return
		}
	}
}

func getFuturePrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceFutureWSURL); err != nil {
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Error read message:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
			continue
		}

		if price, ok := data["p"]; ok {
			formattedPrice := formatNumber(price)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Future price | %s |  %s", symbol, formattedPrice))
			bot.Send(msg)
			return
		}
	}
}

func getFundingRate(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceFutureWSURL); err != nil {
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Error read message:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
			continue
		}

		if fundingRate, ok := data["r"]; ok {
			formattedFundingRate := formatNumber(fundingRate)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding rate | %s |  %s", symbol, formattedFundingRate))
			bot.Send(msg)
			return
		}
	}
}

func getFundingRateCountdown(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceFutureWSURL); err != nil {
		log.Println("Error connect WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Error subscribe stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Error read message:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Error unmarshal json:", err)
			continue
		}

		if nextFundingTime, ok := data["T"]; ok {
			nextFundingTimeMillis := int64(nextFundingTime.(float64))
			remainingTime := time.Until(time.UnixMilli(nextFundingTimeMillis))
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding rate countdown | %s |  %s", symbol, remainingTime))
			bot.Send(msg)
			return
		}
	}
}
