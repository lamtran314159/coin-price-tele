package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"telegram-bot/bot"
	"telegram-bot/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := config.GetEnv("BOT_TOKEN")
	tgBot, err := bot.InitBot(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start the bot to listen for updates
	go bot.Start(ctx, tgBot)

	// Stop the bot when user presses enter
	log.Println("Bot is running. Press enter to stop.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}
