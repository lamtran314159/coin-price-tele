package bot

import (
	"context"
	"log"
	"telegram-bot/bot/handlers"

	// "time"
	"encoding/json"
	"fmt"
	"net/http"

	// "sync"
	// "bytes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

var commands = []tgbotapi.BotCommand{
	{
		Command:     "start",
		Description: "Authenticate and start using the bot",
	},
	{
		Command:     "scream",
		Description: "Enable screaming mode",
	},
	{
		Command:     "whisper",
		Description: "Disable screaming mode",
	},
	{
		Command:     "menu",
		Description: "Show menu with buttons",
	},
	{
		Command:     "help",
		Description: "Show available commands",
	},
	{
		Command:     "protected",
		Description: "Test to see if user is authenticated",
	},
	{
		Command:     "kline",
		Description: "<symbol> <interval> [limit] [startTime] [endTime]",
	},
	{
		Command:     "price_spot",
		Description: "<symbol>",
	},
	{
		Command:     "price_future",
		Description: "<symbol>",
	},
	{
		Command:     "funding_rate",
		Description: "<symbol>",
	},
	{
		Command:     "funding_rate_countdown",
		Description: "<symbol>",
	},
}

type CoinPriceUpdate struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Threshold float64 `json:"threshold"`
	Lower     bool    `json:"lower"`
	VipRole   int     `json:"vip_role"`
	ChatID    int64   `json:"chatID"`
	Timestamp string  `json:"timestamp"`
}

// Initialize the bot with the token
func InitBot(token string, webhookURL string) (*tgbotapi.BotAPI, error) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false // Set to true if you want to debug interactions
	webhook, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return nil, err
	}
	_, err = bot.Request(webhook)
	if err != nil {
		return nil, err
	}
	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Start")
	return bot, nil
}

// Start listening for updates
func Start(ctx context.Context, bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Pass updates to handler
	go receiveUpdates(ctx, updates)
}

// Start listening update from webhook
func StartWebhook(bot *tgbotapi.BotAPI) {
	//Create the update channel using ListenForWebhook
	updates := bot.ListenForWebhook("/webhook")
	for update := range updates {
		if update.Message != nil {
			handlers.HandleMessage(update.Message, bot)
		} else if update.CallbackQuery != nil {
			handlers.HandleButton(update.CallbackQuery, bot)
		}
	}
}

// Receive updates and pass them to handlers
func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			if update.Message != nil {
				handlers.HandleMessage(update.Message, bot)
			} else if update.CallbackQuery != nil {
				handlers.HandleButton(update.CallbackQuery, bot)
			}
		}
	}
}

func PriceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//? nhan lenh post -> gui cho user
	//? print user
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var update CoinPriceUpdate
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Process the received data
	fmt.Printf("Received price update: Coin: %s, Price: %.2f, Timestamp: %s\n", update.Symbol, update.Threshold, update.Timestamp)
	// Sử dụng WaitGroup để quản lý các goroutine
	direction := "below"
	if !update.Lower {
		direction = "above"
	}
	mess := fmt.Sprintf("Price alert: Coin: %s is %s threshold: %.2f\n Current price: %.2f", update.Symbol, direction, update.Threshold, update.Price)
	go handlers.SendMessageToUser(bot, update.ChatID, mess)

	// Respond to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price update received"))
}

//	demo payload {
// 	"symbol": "BTC",
// 	"price": 65000,
// 	"threshold": 60000,
// 	"lower": false,
// 	"vip_role": 1,
// 	"chatID": 6989009560,
// 	"timestamp": "2024-01-01T00:00:00Z"
// }
