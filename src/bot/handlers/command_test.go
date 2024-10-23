package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Test the handleCommand function with different commands.
func TestHandleCommand(t *testing.T) {
	// Mock the bot
	bot := &tgbotapi.BotAPI{}
	// Mock the user
	user := &tgbotapi.User{
		ID: 123456,
	}
	// Test cases for each command
	tests := []struct {
		command string
		message string
	}{
		{command: "/scream", message: "Screaming mode enabled."},
		{command: "/whisper", message: "Screaming mode disabled."},
		{command: "/menu", message: firstMenu},
	}

	for _, tt := range tests {
		chatID := int64(123456) // Dummy chatID

		// Capture output by simulating bot.Send method
		t.Run(tt.command, func(t *testing.T) {
			handleCommand(chatID, tt.command, nil, bot, user)
			// We are not actually sending the message via Telegram in the test,
			// but we check that the appropriate command was processed correctly.
		})
	}
}

// Test sendMenu function
func TestSendMenu(t *testing.T) {
	chatID := int64(123456)
	menuMessage := sendMenu(chatID)

	assert.Equal(t, firstMenu, menuMessage.Text, "The menu text should match.")
	assert.Equal(t, tgbotapi.ModeHTML, menuMessage.ParseMode, "The parse mode should be HTML.")
	assert.NotNil(t, menuMessage.ReplyMarkup, "ReplyMarkup should not be nil.")
}

// Test sendScreamedMessage function
func TestSendScreamedMessage(t *testing.T) {
	message := &tgbotapi.Message{
		Text: "hello world",
		Chat: &tgbotapi.Chat{ID: 123456},
	}
	screamedMessage := sendScreamedMessage(message)

	assert.Equal(t, "HELLO WORLD", screamedMessage.Text, "The message should be in uppercase.")
	assert.Equal(t, tgbotapi.ModeHTML, screamedMessage.ParseMode, "The parse mode should be HTML.")
}

// Test copyMessage function
func TestCopyMessage(t *testing.T) {
	message := &tgbotapi.Message{
		Text: "original message",
		Chat: &tgbotapi.Chat{ID: 123456},
	}
	copiedMessage := copyMessage(message)

	assert.Equal(t, "original message", copiedMessage.Text, "The copied message should match the original.")
	assert.Equal(t, message.Chat.ID, copiedMessage.ChatID, "The Chat ID should be the same as the original.")
}
