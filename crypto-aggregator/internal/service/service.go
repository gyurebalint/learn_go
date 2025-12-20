package service

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/store"
	"time"
)

type PriceService struct {
	db      store.Storage
	fetcher fetcher.Fetcher
}

func NewPriceService(db store.Storage, f fetcher.Fetcher) *PriceService {
	return &PriceService{
		db:      db,
		fetcher: f,
	}
}

func (s *PriceService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	priceValue, err := s.fetcher.Fetch(ctx, symbol)
	if err != nil {
		return 0, err
	}

	price := store.Price{
		Exchange:  "BINANCE",
		Value:     priceValue,
		Currency:  symbol,
		TimeStamp: time.Now()}
	err = s.db.SavePrice(ctx, []store.Price{price})
	if err != nil {
		return 0, err
	}

	return price.Value, nil
}
