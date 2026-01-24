package main

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/service"
	"crypto-aggregator/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	connString := getDbConnectionString()
	db, err := store.NewPostgres(connString)
	if err != nil {
		log.Fatalf("Failed to initiate database: %v", err)
	}
	defer db.Close()

	redisAddr := getRedisAddr()
	rdb, err := store.NewRedis(redisAddr)
	if err != nil {
		log.Fatalf("failed to initialize redis: %v", err)
	}

	binFetcher := &fetcher.BinanceFetcher{}
	kucoinFetcher := &fetcher.KucoinFetcher{}
	serv := service.NewPriceService(db, []fetcher.Fetcher{binFetcher, kucoinFetcher}, rdb)

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
			fmt.Printf("Error processing %s: %v\n", symbol, err)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server running on port: %s...", port)
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
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
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "admin"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "crypto-aggregator"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, host, dbName)
}

func getRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return "localhost:6379"
	}
	return addr
}
