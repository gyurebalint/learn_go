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
type binanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (f *BinanceFetcher) Fetch(ctx context.Context, symbol string) (Response, error) {
	fetcherResp := Response{"BINANCE", 0}

	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fetcherResp, fmt.Errorf("error fetching the price from Binance: %w", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fetcherResp, fmt.Errorf("binance is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fetcherResp, fmt.Errorf("request was not successful: %d", resp.StatusCode)
	}

	var priceResponse binanceResponse
	err = json.NewDecoder(resp.Body).Decode(&priceResponse)
	if err != nil {
		return fetcherResp, fmt.Errorf("failed to decode binanceResponse: %w", err)
	}

	price, err := strconv.ParseFloat(priceResponse.Price, 64)
	if err != nil {
		return fetcherResp, fmt.Errorf("invalid price format '%s': %w", priceResponse.Price, err)
	}

	fetcherResp.Price = price
	return fetcherResp, nil
}
