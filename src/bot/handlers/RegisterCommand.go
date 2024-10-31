package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token = "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJNSyIsInN1YiI6InRyYW5odXkiLCJwYXNzd29yZCI6ImFpIGNobyBjb2kgbeG6rXQga2jhuql1IiwiZXhwIjoxNzMwMzM2MTUwfQ.1N--Ur3cno-j_pFIs9hxc3q6BWLxL0JYWCd9plPx7qE"

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

	req.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
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

func RegisterPriceDifferenceAndFundingRate(ID int64, symbol string, threshold float64, is_lower bool, Type string, bot *tgbotapi.BotAPI) error {
	url := fmt.Sprintf("https://hcmutssps.id.vn/api/vip2/create?triggerType=%s", Type)
	fmt.Println("Type:", Type)
	method := "POST"

	condition := ">="
	if is_lower {
		condition = "<"
	}

	var payload io.Reader
	if Type == "price-difference" {
		payload = strings.NewReader(fmt.Sprintf(`{
	"symbol": "%s",
	"condition": "%s",
	"priceDifference": %f
	}`, symbol, condition, threshold))
		fmt.Println("payload:", payload)
	} else if Type == "funding-rate" {
		payload = strings.NewReader(fmt.Sprintf(`{
	"symbol": "%s",
	"condition": "%s",
	"fundingRate": %f
	}`, symbol, condition, threshold))
		fmt.Println("payload:", payload)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cookie", fmt.Sprintf("token=%s", token))

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
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s of %s below %f threshold successfully!", Type, symbol, threshold)))
		} else {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s of %s above %f threshold successfully!", Type, symbol, threshold)))
		}
	}
	return nil
}

func DeleteTrigger(ID int64, bot *tgbotapi.BotAPI, symbol string, price_type string) {
	url := fmt.Sprintf("https://hcmutssps.id.vn/api/vip2/delete/%s?triggerType=%s", symbol, price_type)
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	bot.Send(tgbotapi.NewMessage(ID, string(body)))
}
