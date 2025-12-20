package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type KucoinFetcher struct {
}
type kucoinResponse struct {
	Code string `json:"code"`
	Data data   `json:"data"`
}

type data struct {
	Price string
	Size  string
}

func (f *KucoinFetcher) Fetch(ctx context.Context, symbol string) (Response, error) {
	fetcherResp := Response{"KUCOIN", 0}

	url := fmt.Sprintf("https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=%s-USDT", symbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fetcherResp, fmt.Errorf("error fetching the price from kucoin: %w", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fetcherResp, fmt.Errorf("kucoin is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fetcherResp, fmt.Errorf("request was not successful: %d", resp.StatusCode)
	}

	var priceResponse kucoinResponse
	err = json.NewDecoder(resp.Body).Decode(&priceResponse)
	if err != nil {
		return fetcherResp, fmt.Errorf("failed to decode kucoinResponse: %w", err)
	}

	price, err := strconv.ParseFloat(priceResponse.Data.Price, 64)
	if err != nil {
		return fetcherResp, fmt.Errorf("invalid price format '%s': %w", priceResponse.Data.Price, err)
	}

	fetcherResp.Price = price
	return fetcherResp, nil
}
