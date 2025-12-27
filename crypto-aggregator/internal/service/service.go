package service

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/store"
	"fmt"
	"strings"
	"sync"
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
	respChan := make(chan fetcher.Response, len(s.fetchers))
	var waitGroup sync.WaitGroup
	for _, f := range s.fetchers {
		waitGroup.Add(1)

		go func(exchange fetcher.Fetcher) {
			defer waitGroup.Done()
			resp, err := exchange.Fetch(ctx, symbol)
			if err != nil {
				return
			}
			fmt.Printf("Fetching from: %s\n", resp.Exchange)
			respChan <- resp
		}(f)
	}
	waitGroup.Wait()
	close(respChan)

	var exchanges []string
	var prices float64
	var validFetcherCount int
	for resp := range respChan {
		if resp.Price == 0 {
			continue
		}

		validFetcherCount += 1
		exchanges = append(exchanges, resp.Exchange)
		prices += resp.Price
	}
	if validFetcherCount == 0 {
		return 0, fmt.Errorf("no valid exchange was available")
	}

	avg := prices / (float64(validFetcherCount))
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

//func (s *PriceService) GetPrice(ctx context.Context, symbol string) (float64, error) {
//	var exchanges []string
//	var prices float64
//	var validFetcherCount int
//
//	for i := range s.fetchers {
//		resp, err := s.fetchers[i].Fetch(ctx, symbol)
//		if err != nil {
//			continue
//		}
//		if resp.Price == 0 {
//			continue
//		}
//
//		validFetcherCount += 1
//		exchanges = append(exchanges, resp.Exchange)
//		prices += resp.Price
//	}
//
//	if validFetcherCount == 0 {
//		return 0, fmt.Errorf("no valid exchange was available")
//	}
//	avg := prices / (float64(validFetcherCount))
//	price := store.Price{
//		Exchange:  strings.Join(exchanges, ","),
//		Value:     avg,
//		Currency:  symbol,
//		TimeStamp: time.Now()}
//	err := s.db.SavePrice(ctx, []store.Price{price})
//	if err != nil {
//		return 0, err
//	}
//
//	return price.Value, nil
//}
