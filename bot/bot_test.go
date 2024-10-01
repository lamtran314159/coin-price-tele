package bot

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestInitBot(t *testing.T) {
	// Test with invalid token
	invalidToken := "invalid_token"
	_, err := InitBot(invalidToken)
	if err == nil {
		t.Error("InitBot() with invalid token should return an error")
	}

	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	validToken := os.Getenv("BOT_TOKEN")
	bot, err := InitBot(validToken)
	if err != nil {
		t.Errorf("InitBot() with valid token returned an error: %v", err)
	}
	if bot == nil {
		t.Error("InitBot() with valid token returned nil bot")
	}
}
