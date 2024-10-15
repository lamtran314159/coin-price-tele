package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}

func AuthenticateUser(telegramId int64) (string, error) {
	url := "http://localhost:3000/auth"

	// Create the request body
	body := map[string]interface{}{
		"telegramId": fmt.Sprintf("%d", telegramId),
	}
	jsonBody, _ := json.Marshal(body)

	// Send the POST request to the mock backend
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("access denied")
	}

	// Decode the response
	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.AccessToken, nil
}

func ValidateToken(token string) (string, error) {
	url := "http://localhost:3000/protected"

	// Add the Authorization header
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	return "", fmt.Errorf("invalid token")
}
