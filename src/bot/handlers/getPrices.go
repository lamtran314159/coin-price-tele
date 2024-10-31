package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	APIBaseURL_Spot_Price    = "https://hcmutssps.id.vn/api/get-spot-price"
	APIBaseURL_Futures_Price = "https://hcmutssps.id.vn/api/get-future-price"
	APIBaseURL_Funding_Rate  = "https://hcmutssps.id.vn/api/get-funding-rate"
	CookieToken              = "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJNSyIsInN1YiI6InRyYW5odXkiLCJwYXNzd29yZCI6ImFpIGNobyBjb2kgbeG6rXQga2jhuql1IiwiZXhwIjoxNzMwMzYyNDI5fQ.R14QMImoy3rRLQus4S9QglMnmdh49BU3Su5v3-rWI4o"
)

type SpotPriceResponse struct {
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
	Symbol    string `json:"symbol"`
}

type FuturesPriceResponse struct {
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
	Symbol    string `json:"symbol"`
}

type FundingRateResponse struct {
	Symbol                   string `json:"symbol"`
	FundingRate              string `json:"fundingRate"`
	FundingRateCountdown     string `json:"fundingCountdown"`
	EventTime                string `json:"eventTime"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int    `json:"fundingIntervalHours"`
}
type PriceInfoSpot struct {
	Symbol    string `json:"Symbol"`
	EventTime string `json:"Event time"`
	SpotPrice string `json:"Spot price"`
}

type PriceInfoFutures struct {
	Symbol    string `json:"Symbol"`
	EventTime string `json:"Event time"`
	SpotPrice string `json:"Futures price"`
}

type PriceInfoFundingRate struct {
	Symbol                   string `json:"Symbol"`
	EventTime                string `json:"Event time"`
	FundingRate              string `json:"Funding rate"`
	FundingRateCountdown     string `json:"Time until next funding"`
	FundingRateIntervalHours string `json:"Funding rate interval"`
}

func formatPrice(input string) string {

	parts := strings.Split(input, ".")

	// Xử lý phần nguyên
	intPart := parts[0]
	n := len(intPart)
	if n <= 3 {
		return input
	}

	var result strings.Builder
	offset := n % 3
	if offset > 0 {
		result.WriteString(intPart[:offset])
		if n > 3 {
			result.WriteString(",")
		}
	}
	for i := offset; i < n; i += 3 {
		result.WriteString(intPart[i : i+3])
		if i+3 < n {
			result.WriteString(",")
		}
	}

	if len(parts) > 1 {
		result.WriteString(".")
		result.WriteString(parts[1])
	}

	return result.String()
}

func intToString(n int) string {
	return strconv.Itoa(n)
}

func FormatPrice1(a string) string {
	// Lặp từ cuối chuỗi và loại bỏ các số 0 ở cuối
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] != '0' {
			if a[i] == '.' {
				return a + "00"
			}
			return a
		}
		// Nếu ký tự cuối là '0', loại bỏ ký tự đó
		a = a[:i]
	}

	return a
}

func GetSpotPriceStream(chatID int64, symbol string, bot *tgbotapi.BotAPI) {

	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure context is canceled when done

	// Create the request URL
	reqUrl := fmt.Sprintf("%s?symbols=%s", APIBaseURL_Spot_Price, symbol)
	//log.Printf("API URL: %s", reqUrl)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return
	}

	// Set necessary headers for the request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", CookieToken))

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch spot price: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received status code %d", resp.StatusCode)
		if resp.StatusCode == 500 {
			msg := tgbotapi.NewMessage(chatID, "You need to authenticate before executing this command.")
			bot.Send(msg)
		}
		return
	}

	// Read data from the stream
	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		// Remove the "data:" prefix from the line
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Decode JSON
		var spotPriceResponse SpotPriceResponse
		err := json.Unmarshal([]byte(line), &spotPriceResponse)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue // Skip the error and continue reading the next data
		}
		pricestr := FormatPrice1(spotPriceResponse.Price)
		//log.Printf("Price: %s", pricestr)
		// price, err := strconv.ParseFloat(pricestr, 64)
		// if err != nil {
		// 	log.Printf("Error converting price to float: %v", err)
		// 	continue
		// }

		// Send decoded data to Telegram user and exit
		if strings.EqualFold(spotPriceResponse.Symbol, symbol) {
			formattedPrice := formatPrice(pricestr)

			priceInfo := PriceInfoSpot{
				Symbol:    spotPriceResponse.Symbol,
				EventTime: spotPriceResponse.EventTime,
				SpotPrice: formattedPrice,
			}
			// Convert the object to a JSON string
			jsonData, err := json.MarshalIndent(priceInfo, "", "    ")
			if err != nil {
				log.Printf("Error creating JSON: %v", err)
				//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating JSON: %v", err)))
				return
			}

			// Use HTML to display the JSON string
			msg := tgbotapi.NewMessage(chatID, "<pre>"+string(jsonData)+"</pre>")
			msg.ParseMode = "HTML"
			bot.Send(msg)

			cancel()
			return // Exit immediately after sending the first data
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
		//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol is not available. Please provide a valid symbol.")))
	}
}

func GetFuturesPriceStream(chatID int64, symbol string, bot *tgbotapi.BotAPI) {

	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure context is canceled when done

	// Create the request URL
	reqUrl := fmt.Sprintf("%s?symbols=%s", APIBaseURL_Futures_Price, symbol)
	log.Printf("API URL: %s", reqUrl)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return
	}

	// Set necessary headers for the request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", CookieToken))

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch futures price: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received status code %d", resp.StatusCode)
		if resp.StatusCode == 500 {
			msg := tgbotapi.NewMessage(chatID, "You need to authenticate before executing this command.")
			bot.Send(msg)
		}
		return
	}

	// Read data from the stream
	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		// Remove the "data:" prefix from the line
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Decode JSON
		var futuresPriceResponse FuturesPriceResponse
		err := json.Unmarshal([]byte(line), &futuresPriceResponse)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue // Skip the error and continue reading the next data
		}
		pricestr := FormatPrice1(futuresPriceResponse.Price)

		// price, err := strconv.ParseFloat(pricestr, 64)
		// if err != nil {
		// 	log.Printf("Error converting price to float: %v", err)
		// 	continue
		// }

		// Send decoded data to Telegram user and exit
		if strings.EqualFold(futuresPriceResponse.Symbol, symbol) {
			formattedPrice := formatPrice(pricestr)

			priceInfo := PriceInfoFutures{
				Symbol:    futuresPriceResponse.Symbol,
				EventTime: futuresPriceResponse.EventTime,
				SpotPrice: formattedPrice,
			}
			// Convert the object to a JSON string
			jsonData, err := json.MarshalIndent(priceInfo, "", "    ")
			if err != nil {
				log.Printf("Error creating JSON: %v", err)
				//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating JSON: %v", err)))
				return
			}

			// Use HTML to display the JSON string
			msg := tgbotapi.NewMessage(chatID, "<pre>"+string(jsonData)+"</pre>")
			msg.ParseMode = "HTML"
			bot.Send(msg)

			cancel()
			return // Exit immediately after sending the first data
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
		//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol is not available. Please provide a valid symbol.")))
	}
}

func GetFundingRateStream(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure context is canceled when done

	// Create the request URL
	reqUrl := fmt.Sprintf("%s?symbols=%s", APIBaseURL_Funding_Rate, symbol)
	log.Printf("API URL: %s", reqUrl)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return
	}

	// Set necessary headers for the request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cookie", fmt.Sprintf("token=%s", CookieToken))

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch funding rate: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received status code %d", resp.StatusCode)
		if resp.StatusCode == 500 {
			msg := tgbotapi.NewMessage(chatID, "You need to authenticate before executing this command.")
			bot.Send(msg)
		}
		return
	}

	// Read data from the stream
	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		// Remove the "data:" prefix from the line
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Decode JSON
		var fundingRateResponse FundingRateResponse
		err := json.Unmarshal([]byte(line), &fundingRateResponse)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue // Skip the error and continue reading the next data
		}
		fundingstr := FormatPrice1(fundingRateResponse.FundingRate)
		// Send decoded data to Telegram user and exit
		if strings.EqualFold(fundingRateResponse.Symbol, symbol) {
			fundingRateInterval := intToString(fundingRateResponse.FundingIntervalHours)
			priceInfo := PriceInfoFundingRate{
				Symbol:                   fundingRateResponse.Symbol,
				EventTime:                fundingRateResponse.EventTime,
				FundingRate:              fundingstr,
				FundingRateCountdown:     fundingRateResponse.FundingRateCountdown,
				FundingRateIntervalHours: fundingRateInterval + " hours",
			}

			// Convert the object to a JSON string
			jsonData, err := json.MarshalIndent(priceInfo, "", "    ")
			if err != nil {
				log.Printf("Error creating JSON: %v", err)
				//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating JSON: %v", err)))
				return
			}

			// Use HTML to display the JSON string
			msg := tgbotapi.NewMessage(chatID, "<pre>"+string(jsonData)+"</pre>")
			msg.ParseMode = "HTML"
			bot.Send(msg)

			cancel()
			return // Exit immediately after sending the first data
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
		//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol is not available. Please provide a valid symbol.")))
	}
}
