package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(connString string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		println("cannot connect to db", err.Error())
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping faild: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS prices (
		id SERIAL PRIMARY KEY,
		exchange TEXT NOT NULL,
		currency TEXT NOT NULL,
		value DOUBLE PRECISION NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL
	);`

	_, err = pool.Exec(context.Background(), query)
	if err != nil {
		println("cannot create prices table", err.Error())
	}
	return &Postgres{pool: pool}, nil
}

func (db *Postgres) SavePrice(ctx context.Context, prices []Price) error {
	query := `INSERT INTO prices (exchange, currency, value,timestamp) VALUES ($1, $2, $3, $4)`

	for _, price := range prices {
		_, err := db.pool.Exec(ctx, query, price.Exchange, price.Currency, price.Value, price.TimeStamp)
		if err != nil {
			fmt.Printf("cannot save price %v: %v", price, err)
		}
	}
	return nil
}

func (db *Postgres) GetLatestPrice(ctx context.Context, currency string) (Price, error) {
	query := `SELECT exchange, currency, value, timestamp
    FROM prices 
    WHERE currency = $1
    ORDER BY timestamp DESC LIMIT 1`

	var price Price
	var err = db.pool.QueryRow(ctx, query, currency).Scan(
		&price.Exchange,
		&price.Currency,
		&price.Value,
		&price.TimeStamp)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Price{}, fmt.Errorf("no price found for %s", currency)
		}
		return Price{}, fmt.Errorf("query failed: %w", err)
	}

	return price, nil
}

func (db *Postgres) Close() {
	db.pool.Close()
}
