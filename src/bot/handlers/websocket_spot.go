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

// Global variables to store prices for multiple symbols
var SpotPrices = make(map[string]float64)
var mu sync.Mutex

// Fetch the list of available symbols from Binance API
func FetchBinanceSymbols_Spot() ([]string, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
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

// Start WebSocket to listen for a single crypto price
func StartWebSocket_Spot(symbol string) {
	ctx := context.Background()
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", strings.ToLower(symbol))

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
		//log.Printf("Failed to connect to WebSocket for %s: %v", symbol, err)
		return // Skip to the next symbol if the connection fails
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Printf("Failed to read message for %s: %v", symbol, err)
			break // Break the loop and exit if reading fails
		}

		var result struct {
			Symbol string `json:"s"`
			Price  string `json:"p"`
		}

		if err := json.Unmarshal(message, &result); err != nil {
			log.Printf("Failed to parse message for %s: %v", symbol, err)
			continue
		}

		price, err := strconv.ParseFloat(result.Price, 64)
		if err != nil {
			log.Printf("Failed to convert price to float64 for %s: %v", symbol, err)
			continue
		}
		// Update price for the respective symbol
		mu.Lock()
		//log.Printf("Price for %s: %.4f", symbol, price)
		SpotPrices[result.Symbol] = price
		mu.Unlock()
	}
}

// Get the price of a cryptocurrency symbol
func GetSpotPrice(symbol string) (float64, bool) {
	mu.Lock()
	defer mu.Unlock()
	price, exists := SpotPrices[symbol]
	if !exists {
		return 0, false
	}
	return price, exists
}

func FetchandStartWebSocket_Spot() {
	symbols, err := FetchBinanceSymbols_Spot()
	if err != nil {
		log.Fatalf("Failed to fetch Binance symbols: %v", err)
	}
	// Start WebSocket to get prices for all symbols
	for _, symbol := range symbols {
		//log.Printf("Starting WebSocket for symbol: %s", symbol)
		go StartWebSocket_Spot(symbol)
	}
}
