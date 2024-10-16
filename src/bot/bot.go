package bot

import (
	"context"
	"log"
	"telegram-bot/bot/handlers"

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
	log.Printf("here")
	updates := bot.ListenForWebhook("/webhook")
	if updates == nil{
		log.Printf("here")
	}
	for update := range updates{
		if update.Message != nil {
			log.Print("got updates")
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
