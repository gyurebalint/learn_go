package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log_ingester/internal/models"
	"log_ingester/internal/storage"
	_ "modernc.org/sqlite"
	"sync"
	"time"
)

func main() {
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
	jsonString := `[
		{
			"person_name":"John Doe", "person_age":42
		},
		{
			"person_name":"Jane Doe", "person_age":43
		},
		{
			"person_name":"Jerry Doe", "person_age":44
		},
		{
			"person_name":"Jermaine Doe", "person_age":45
		},
		{
			"person_name":"Jeremiah Doe", "person_age":46
		}
]`

	var people []models.Person
	err = json.Unmarshal([]byte(jsonString), &people)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
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
	for _, person := range people {
		ch <- person
	}
	close(ch)
	workerGroup.Wait()
	close(errCh)
	errGroup.Wait()

	fmt.Println("Refactor complete")
}
