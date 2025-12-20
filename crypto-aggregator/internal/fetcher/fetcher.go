package fetcher

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, symbol string) (float64, error)
}
