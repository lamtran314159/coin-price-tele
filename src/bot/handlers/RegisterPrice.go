package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ErrorResponse struct {
	AlertID string `json:"alert_id"`
	Message string `json:"message"`
}

func RegisterPriceThreshold(ID int64, symbol string, threshold float64, is_lower bool, price_type string, bot *tgbotapi.BotAPI) error {
	url := fmt.Sprintf("https://hcmutssps.id.vn/api/vip2/create?triggerType=%s", price_type)
	fmt.Println("price_type:", price_type)
	method := "POST"

	condition := ">="
	if is_lower {
		condition = "<"
	}

	payload := strings.NewReader(fmt.Sprintf(`{
	  "symbol": "%s",
	  "price": %f,
	  "condition": "%s"
	}`, symbol, threshold, condition))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	//token có thể thay đổi
	req.Header.Add("Cookie", "token=eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJNSyIsInN1YiI6InRyYW5odXkiLCJwYXNzd29yZCI6ImFpIGNobyBjb2kgbeG6rXQga2jhuql1IiwiZXhwIjoxNzMwMTgyOTg2fQ.M0HlUeoRE5WuA59oIDWu9uQ32U1rIQIGd9AQiTlEbX4")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}
	//bot.Send(tgbotapi.NewMessage(ID, errorResponse.Message))
	if errorResponse.AlertID != "" {
		if condition == "<" {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s price of %s below %f threshold successfully!", price_type, symbol, threshold)))
		} else {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s price of %s above %f threshold successfully!", price_type, symbol, threshold)))
		}
	}
	return nil
}

func RegisterPriceDifference() {}

func RegisterFundingRateChange() {}
