package bot

import (
	"context"
	"log"
	"telegram-bot/bot/handlers"

	// "time"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	{
		Command:     "alert_price_with_threshold",
		Description: "<spot/future> <lower/above> <symbol> <threshold>",
	},
	{
		Command:     "price_differece",
		Description: "<threshold>",
	},
	{
		Command:     "funding_rate_change",
		Description: "<threshold>",
	},
}


// send from BE
type CoinPriceUpdate struct {
	Symbol      string  `json:"symbol"`
	Spotprice   float64 `json:"spot_price"`
	Futureprice float64 `json:"future_price"`
	Pricediff   float64 `json:"price_diff"`
	Fundingrate float64 `json:"fundingrate"`
	Threshold   float64 `json:"threshold"`
	Condition   string  `json:"condition"`
	ChatID      string  `json:"chatID"`
	Timestamp   string  `json:"timestamp"`
	Triggertype string  `json:"triggerType"` //spot, price-difference, funding-rate, future

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

	if (update.Condition == ">=" || update.Condition == ">") {
		direction = "above"
	}
	chatID, err := strconv.ParseInt(update.ChatID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	var mess string
	if update.Triggertype == "spot" {
		mess = fmt.Sprintf("Price alert: Coin: %s is %s threshold: %.2f\n Current spot price: %.2f\n Trigger Type: %s",

							update.Symbol, direction, update.Threshold, update.Spotprice, update.Triggertype)
	}
	if update.Triggertype == "price-difference" {
		mess = fmt.Sprintf("Price alert: Coin: %s is %s Price-diff: %.2f\n Current spot price: %.2f, Current future price: %.2f\n Trigger Type: %s",
							update.Symbol, direction, update.Pricediff, update.Spotprice, update.Futureprice, update.Triggertype)
	}
	go handlers.SendMessageToUser(bot, chatID, mess)

	// Respond to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price update received"))
}

// demo payload{
//     "symbol": "BTC",
//     "spot_price": 65000,
//     "future_price": 64000,
//     "threshold": 60000,
//     "condition" : ">="
//     "chatID": "6989009560",
//     "timestamp": "2024-01-01T00:00:00Z",
//     "triggerType": "price-difference"
// }
