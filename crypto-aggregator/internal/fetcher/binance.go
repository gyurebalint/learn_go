package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type BinanceFetcher struct {
}
type response struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (f *BinanceFetcher) Fetch(ctx context.Context, symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error fetching the price from Binance: %w", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("binance is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("request was not successful: %d", resp.StatusCode)
	}

	var binanceResponse response
	err = json.NewDecoder(resp.Body).Decode(&binanceResponse)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	price, err := strconv.ParseFloat(binanceResponse.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price format '%s': %w", binanceResponse.Price, err)
	}

	return price, nil
}
