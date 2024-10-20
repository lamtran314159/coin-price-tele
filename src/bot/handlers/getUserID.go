package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "io"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)   

// Struct to send chatID to the backend
type ChatIDPayload struct {
    ChatID int64 `json:"chat_id"`
    // threshold float64 `json:"threshold"`
}
type chatIDResponse struct {
     ChatIDs []int64 `json:"chat_ids"`
    // threshold float64 `json:"threshold"`
}

// Function to store chatID in the backend
func StoreChatID(chatID int64) error {
    url := "https://your-backend-api.com/store-chat-id" // Replace with your backend API URL

    // Create payload with chatID
    payload := ChatIDPayload{
        ChatID: chatID,
    }

    // Convert payload to JSON
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("error marshalling chatID payload: %v", err)
    }

    // Send POST request to backend API
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return fmt.Errorf("error sending request to backend: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("error response from backend: %s", resp.Status)
    }

    log.Printf("Stored chatID %d successfully!", chatID)
    return nil
}

// Function to retrieve chatIDs from the backend
func GetChatIDs() ([]int64, error) {
    url := "http://localhost:3000/auth"  // Replace with your backend API URL

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
    var chatIDResponse chatIDResponse
    err = json.Unmarshal(body, &chatIDResponse)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling response: %v", err)
    }

    return chatIDResponse.ChatIDs, nil
}

func NotifyUsers(bot *tgbotapi.BotAPI) {
    chatIDs, err := GetChatIDs()
    if err != nil {
        log.Fatalf("Error retrieving chatIDs: %v", err)
    }

    for _, chatID := range chatIDs {
        msg := tgbotapi.NewMessage(chatID, "Your price alert has been triggered!")
        bot.Send(msg)
    }
}