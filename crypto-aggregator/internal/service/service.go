package service

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/store"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PriceService struct {
	db       store.Storage
	fetchers []fetcher.Fetcher
	cache    *store.RedisClient
}

func NewPriceService(db store.Storage, f []fetcher.Fetcher, cache *store.RedisClient) *PriceService {
	return &PriceService{
		db:       db,
		fetchers: f,
		cache:    cache,
	}
}

func (s *PriceService) GetPrice(ctx context.Context, symbol string) (float64, error) {
	cacheKey := fmt.Sprintf("price:%s", symbol)
	cachedVal, err := s.cache.Get(ctx, cacheKey)

	if err == nil {
		var price float64
		fmt.Printf("Cache HIT for %s\n", symbol)
		price, err = strconv.ParseFloat(cachedVal, 64)
		if err != nil {
			return 0.0, fmt.Errorf("failed to parse return cached value to price float64, err: %s", err)
		}
		return price, nil
	}

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
	err = s.db.SavePrice(ctx, []store.Price{price})
	if err != nil {
		return 0, err
	}

	err = s.cache.Set(ctx, cacheKey, price.Value, 60*time.Second)
	if err != nil {
		fmt.Printf("⚠️ Failed to write to Redis: %v\n", err)
	}
	fmt.Printf("💾 Saved to Redis: %s = %f\n", cacheKey, price.Value)

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
