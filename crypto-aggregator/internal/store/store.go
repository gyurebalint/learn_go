package store

import (
	"context"
	"time"
)

type Price struct {
	Exchange  string
	Currency  string //"TAO", "BTC"
	Value     float64
	TimeStamp time.Time
}

type Storage interface {
	SavePrice(ctx context.Context, prices []Price) error
	GetLatestPrice(ctx context.Context, currency string) (Price, error)
	Close()
}
