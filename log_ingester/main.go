package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log_ingester/internal/models"
	"log_ingester/internal/storage"
	_ "modernc.org/sqlite"
	"os"
	"sync"
	"time"
)

func main() {
	inputFile := flag.String("file", "users.json", "Input File")
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Println("file could not be opened", err)
		return
	}

	//Init DB
	dest, err := storage.NewSqliteDest("data.db")
	if err != nil {
		log.Println("Failed to init DB", err)
		return
	}

	defer dest.Close()

	//Init ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	ch := make(chan models.Person)
	errCh := make(chan error)

	var workerGroup sync.WaitGroup
	var errGroup sync.WaitGroup

	errGroup.Add(1)
	go func() {
		defer errGroup.Done()
		for err := range errCh {
			fmt.Println("DB ERROR", err.Error())
		}
	}()

	//listener
	for range 3 {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			for p := range ch {
				if err := dest.Save(ctx, p); err != nil {
					errCh <- err
				} else {
					fmt.Printf("Inserted person: %s\n", p.Name)
				}
			}
		}()
	}

	//producer
	//reader := strings.NewReader(jsonString)
	dec := json.NewDecoder(file)
	_, err = dec.Token()
	if err != nil {
		log.Println("error reading first token", err)
	}

	for dec.More() {
		var p models.Person
		if err := dec.Decode(&p); err != nil {
			log.Println("error decoding person from data", err)
		}
		ch <- p
	}
	close(ch)
	workerGroup.Wait()
	close(errCh)
	errGroup.Wait()

	fmt.Println("Refactor complete")
}
