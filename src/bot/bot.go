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
}

type CoinPriceUpdate struct {
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
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
	fmt.Printf("Received price update: Coin: %s, Price: %.2f, Timestamp: %s\n", update.Coin, update.Price, update.Timestamp)
	// Sử dụng WaitGroup để quản lý các goroutine
	go handlers.NotifyUsers(bot)

	// Respond to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price update received"))
}

type DemoPayload struct {
	TelegramUserID int64 `json:"telegram_user_id"`
}

// Function to response to /backend endpoint then send message to user
func BackendHandler(w http.ResponseWriter, r *http.Request) {
	var demoPayload DemoPayload
	if err := json.NewDecoder(r.Body).Decode(&demoPayload); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	handlers.SendMessageToUser(bot, demoPayload.TelegramUserID, "Hello, World!")
	w.WriteHeader(http.StatusOK)
}
