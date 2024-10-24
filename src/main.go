package main

import (
	"log"
	"net/http"
	"os"
	"telegram-bot/bot"
	"telegram-bot/config"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	botToken := config.GetEnv("BOT_TOKEN")
	webhookURL := config.GetEnv("WEBHOOK_URL")
	tgBot, err := bot.InitBot(botToken, webhookURL)
	if err != nil {
		log.Panic(err)
	}

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8443"
	}
	go http.HandleFunc("/backend", bot.BackendHandler)

	go http.ListenAndServe(":"+port, nil)
	log.Printf("Bot is listening on port %s...\n", port)
	bot.StartWebhook(tgBot)
}
