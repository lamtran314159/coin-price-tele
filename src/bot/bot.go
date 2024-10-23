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

//Start listening update from webhook
func StartWebhook(bot *tgbotapi.BotAPI){
	//Create the update channel using ListenForWebhook
	updates := bot.ListenForWebhook("/webhook")
	for update := range updates{
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

// func Handlebackend(){



// 	// Giả sử đây là danh sách các chatID đọc từ backend
// 	// sendGetRequest()
//     chatIDs := []int64{123456789, 987654321, 1122334455}

//     // Nội dung tin nhắn
//     message := "Hello from your Telegram bot!"

//     // Sử dụng WaitGroup để quản lý các goroutine
//     var wg sync.WaitGroup

//     // Gửi tin nhắn cho từng chatID
//     for _, chatID := range chatIDs {
//         wg.Add(1) // Tăng số lượng công việc đang chờ
//         go sendMessage(chatID, message, &wg) // Khởi chạy goroutine để gửi tin nhắn
//     }

//     // Đợi tất cả goroutine hoàn thành
//     wg.Wait()

//     fmt.Println("All messages sent!")
// }
// const (
//     btcThreshold = 65000.0 // Set your threshold here
//     checkInterval = 1 * time.Minute // Check every minute
// )

// // PriceResponse represents the response from the Binance API
// type PriceResponse struct {
//     Symbol string `json:"symbol"`
//     Price  string `json:"price"`
// }

// // Function to fetch the current BTC price from Binance
// func fetchBTCPrice(symbol string) (float64, error) {
// 	//Symbol and threadhold set by user
// 	log.Printf("https://api.binance.com/api/v3/ticker/price?symbol="+symbol)
//     resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol="+symbol)
//     if err != nil {
//         return 0, err
//     }
//     defer resp.Body.Close()

//     var priceResponse PriceResponse
//     if err := json.NewDecoder(resp.Body).Decode(&priceResponse); err != nil {
//         return 0, err
//     }

//     // Convert price to float64
//     var price float64
//     if _, err := fmt.Sscanf(priceResponse.Price, "%f", &price); err != nil {
//         return 0, err
//     }

//     return price, nil
// }

// // Function to monitor BTC price
// func MonitorBTCPrice(bot *tgbotapi.BotAPI, chatID int64, symbol string) {
//     for {
//         price, err := fetchBTCPrice(symbol)
//         if err != nil {
//             log.Println("Error fetching %s price:", symbol, err)
//             continue
//         }

//         log.Printf("Current %s price: %.2f USDT", symbol, price)

//         if price > btcThreshold {
//             msg := tgbotapi.NewMessage(chatID, "🚨 Alert: price has exceeded 65,000 USDT! Current price: "+fmt.Sprintf("%.2f", price))
//             bot.Send(msg)
//         }

//         time.Sleep(checkInterval)
//     }
// }
