package fetcher

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, symbol string) (Response, error)
}

type Response struct {
	Exchange string
	Price    float64
}
