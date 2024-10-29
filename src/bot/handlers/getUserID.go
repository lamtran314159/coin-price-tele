package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Payload struct {
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	VipRole   int       `json:"vip_role"`
	IpList    []string  `json:"ip_list"`
	ChatID    int64     `json:"chatID"`
	Symbol    []string  `json:"symbol"`
	Threshold []float64 `json:"threshold"`
}


var storedChatIDs = make(map[int64]bool)

//!send to BE
type CoinPriceUpdate struct {
	Symbol   string    `json:"symbol"`
	Price float64  `json:"threshold"`
	Condition string `json:"condition"` // >= price, <=, >, < 
	Triggertype string `json:"triggerType"`
}

// Function to store chatID in the backend
func StoreChatID(ID int64) error {
	if storedChatIDs[ID] {
		fmt.Printf("ChatID %d already stored!\n", ID)
		return nil
	}
	storedChatIDs[ID] = true

	url := "https://1662-116-110-41-111.ngrok-free.app/backend" // Replace with your backend API URL

	// Create payload with chatID
	payload := Payload{
		ChatID: ID,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("error marshalling chatID payload: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))

	if err != nil {
		return fmt.Errorf("error sending request to backend: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body: %v", err)
			return fmt.Errorf("error reading response body: %v", err)
		}
		fmt.Printf("backend returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("backend returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}
	fmt.Printf("Stored chatID %d successfully!\n", ID)
	return nil
}

func SendMessageToUser(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

// Function to retrieve chatIDs from the backend
func GetChatIDs() /*([]int64, error)*/ ([]string, error) {
	url := "http://103.205.60.174:8080/admin/getAllUser" // Replace with your backend API URL

	// Send GET request to backend
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request to backend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from backend: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response
	var users []Payload
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// var chatIDs [] int64
	// for _, user := range users{
	//     chatIDs = append(chatIDs, user.ChatID)
	// }

	var chatIDs []string
	for _, user := range users {
		chatIDs = append(chatIDs, user.Username)
	}

	return chatIDs, nil
}

// Function to update chatID, symbol, and threshold in the backend
func UpdateChatIDSymbolThreshold(chatID int64, symbol []string, threshold []float64) error {
	url := "http://103.205.60.174:8080/....."

	// Create payload with chatID, symbol, and threshold
	payload := struct {
		ChatID    int64     `json:"chatID"`
		Symbol    []string  `json:"symbol"`
		Threshold []float64 `json:"threshold"`
	}{
		ChatID:    chatID,
		Symbol:    symbol,
		Threshold: threshold,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	// Send PUT request to backend API
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to backend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from backend: %s", resp.Status)
	}

	log.Printf("Updated chatID %d, symbol %v, and threshold %v successfully!", chatID, symbol, threshold)
	return nil
}

// func NotifyUsers(bot *tgbotapi.BotAPI) {
// 	chatIDs, err := GetChatIDs()
// 	if err != nil {
// 		log.Fatalf("Error retrieving chatIDs: %v", err)
// 	}

// 	for _, chatID := range chatIDs {
// 		msg := tgbotapi.NewMessage(6989009560, chatID)
// 		bot.Send(msg)
// 	}
// }
