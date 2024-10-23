package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	//"time"

	"github.com/coder/websocket"
)

// Global variables to store prices for multiple symbols and funding rates

var FuturesPrices = make(map[string]float64)
var FuturesFundingRates = make(map[string]float64)       // New map for funding rates
var FuturesFundingRateCountdown = make(map[string]int64) // New map for funding rate countdown
var mu1 sync.Mutex

// Fetch the list of available symbols from Binance API (for Spot trading)
func FetchBinanceSymbols_Futures() ([]string, error) {
	resp, err := http.Get("https://fapi.binance.com/fapi/v1/exchangeInfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Symbols []struct {
			Symbol string `json:"symbol"`
		} `json:"symbols"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Accept all symbols without filtering
	var symbols []string
	for _, s := range result.Symbols {
		//log.Printf("Symbol: %s", s.Symbol)
		symbols = append(symbols, s.Symbol)
	}
	return symbols, nil
}

// Start WebSocket to listen for Futures data, including funding rate and countdown
func StartWebSocket_Futures(symbol string) {
	ctx := context.Background()
	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@markPrice@1s", strings.ToLower(symbol))

	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPClient: &http.Client{},
	})
	/*client := &http.Client{
		Timeout: 20 * time.Second, // Set timeout to 20 seconds
	}

	// Attempt to connect to WebSocket
	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPClient: client,
	})*/

	if err != nil {
		//log.Printf("Failed to connect to WebSocket for Futures %s: %v", symbol, err)
		return // Skip to the next symbol if the connection fails
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	// Connection succeeded, start reading messages
	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Printf("Failed to read message for Futures %s: %v", symbol, err)
			break // Break the loop if reading fails
		}

		var result struct {
			Symbol          string `json:"s"` // Symbol
			MarkPrice       string `json:"p"` // Mark Price
			FundingRate     string `json:"r"` // Funding Rate
			NextFundingTime int64  `json:"T"` // Time until next funding rate (in milliseconds)
		}
		if err := json.Unmarshal(message, &result); err != nil {
			log.Printf("Failed to parse message for Futures %s: %v", symbol, err)
			continue
		}

		markPrice, err := strconv.ParseFloat(result.MarkPrice, 64)
		if err != nil {
			log.Printf("Failed to convert Futures price to float64 for %s: %v", symbol, err)
			continue
		}

		fundingRate, err := strconv.ParseFloat(result.FundingRate, 64)
		if err != nil {
			log.Printf("Failed to convert Futures funding rate to float64 for %s: %v", symbol, err)
			continue
		}

		// Update mark price, funding rate, and countdown
		mu1.Lock()
		FuturesPrices[result.Symbol] = markPrice
		FuturesFundingRates[result.Symbol] = fundingRate
		FuturesFundingRateCountdown[result.Symbol] = result.NextFundingTime
		mu1.Unlock()
	}
}

// Get the price of a cryptocurrency symbol
func GetFuturePrice(symbol string) (float64, bool) {
	mu1.Lock()
	defer mu1.Unlock()
	price, exists := FuturesPrices[symbol]
	if !exists {
		return 0, false
	}
	return price, exists
}

// Get the futures funding rate of a cryptocurrency symbol
func GetFuturesFundingRate(symbol string) (float64, bool) {
	mu1.Lock()
	defer mu1.Unlock()
	rate, exists := FuturesFundingRates[symbol]
	if !exists {
		return 0, false
	}
	return rate, exists
}

// Get the funding rate countdown (time until next funding) of a cryptocurrency symbol
func GetFundingRateCountdown(symbol string) (int64, bool) {
	mu1.Lock()
	defer mu1.Unlock()
	countdown, exists := FuturesFundingRateCountdown[symbol]
	if !exists {
		return 0, false
	}
	return countdown, exists
}

// Fetch symbols and start WebSocket for Futures
func FetchAndStartWebSocket_Futures() {
	symbols, err := FetchBinanceSymbols_Futures()
	if err != nil {
		log.Fatalf("Failed to fetch Binance symbols: %v", err)
	}
	// Start WebSocket to get futures data for all symbols
	for _, symbol := range symbols {
		go StartWebSocket_Futures(symbol)
	}
}
