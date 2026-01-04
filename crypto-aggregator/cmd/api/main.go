package main

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/service"
	"crypto-aggregator/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	connString := getDbConnectionString()
	db, err := store.NewPostgres(connString)
	if err != nil {
		panic(err)
	}

	binFetcher := &fetcher.BinanceFetcher{}
	kucoinFetcher := &fetcher.KucoinFetcher{}
	serv := service.NewPriceService(db, []fetcher.Fetcher{binFetcher, kucoinFetcher})

	http.HandleFunc("/price", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET params were:", r.URL.Query())

		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			http.Error(w, "missing 'symbol' query parameter", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		price, err := serv.GetPrice(ctx, symbol)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", symbol, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(map[string]any{
			"symbol": symbol,
			"price":  price,
		}); err != nil {
			fmt.Println("Error writing response:", err)
		}
	})

	fmt.Println("Server running on port 3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("could not start server")
	}
}
func getDbConnectionString() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "admin"
	}

	return fmt.Sprintf("postgres://%s:admin@%s:5432/crypto-aggregator", user, host)
}
