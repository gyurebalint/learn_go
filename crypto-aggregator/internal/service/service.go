package service

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/store"
	"strings"
	"time"
)

type PriceService struct {
	db       store.Storage
	fetchers []fetcher.Fetcher
}

func NewPriceService(db store.Storage, f []fetcher.Fetcher) *PriceService {
	return &PriceService{
		db:       db,
		fetchers: f,
	}
}

func (s *PriceService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	var exchanges []string
	var prices float64
	var validFetcherCount int

	for i := range s.fetchers {
		resp, err := s.fetchers[i].Fetch(ctx, symbol)
		if err != nil {
			continue
		}
		if resp.Price == 0 {
			continue
		}

		validFetcherCount += 1
		_ = append(exchanges, resp.Exchange)
		prices += resp.Price
	}

	avg := prices / (float64(len(s.fetchers)))
	price := store.Price{
		Exchange:  strings.Join(exchanges, ","),
		Value:     avg,
		Currency:  symbol,
		TimeStamp: time.Now()}
	err := s.db.SavePrice(ctx, []store.Price{price})
	if err != nil {
		return 0, err
	}

	return price.Value, nil
}
