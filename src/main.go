package main

import (
	// "bufio"
	// "context"
	"log"
	// "os"
	"telegram-bot/bot"
	"telegram-bot/config"
	"net/http"
	"telegram-bot/bot/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := config.GetEnv("BOT_TOKEN")
	webhookURL := config.GetEnv("WEBHOOK_URL")
	tgBot, err := bot.InitBot(botToken, webhookURL)
	if err != nil {
		log.Panic(err)
	}

	// // Create a cancellable context
	// ctx, cancel := context.WithCancel(context.Background())

	// // Start the bot to listen for updates
	// go bot.Start(ctx, tgBot)

	// // Stop the bot when user presses enter
	// log.Println("Bot is running. Press enter to stop.")
	// bufio.NewReader(os.Stdin).ReadBytes('\n')
	// cancel()

	// Start an HTTP server to listen for incoming requests
	port := "8443"
	go http.ListenAndServe(":"+port, nil)
	log.Printf("Bot is listening on port %s...\n",port)
	// go bot.StartWebhook(tgBot)

	// Handle incoming updates from the update channel
	updates := tgBot.ListenForWebhook("/webhook")

	for update := range updates{
		if update.Message != nil {
			log.Print("got updates")
			handlers.HandleMessage(update.Message, tgBot)
		} else if update.CallbackQuery != nil {
			handlers.HandleButton(update.CallbackQuery, tgBot)
		}
	}
}
