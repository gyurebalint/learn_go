package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "modernc.org/sqlite"
	"sync"
	"time"
)

type DataDest interface {
	Save(ctx context.Context, p Person) error
}

type ConsoleDest struct{}

func (ConsoleDest) Save(ctx context.Context, p Person) error {
	select {
	case <-time.After(time.Second * 3):
		fmt.Println("Saved")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

	if p.Name == "John Doe" {
		return errors.New("database connection lost")
	}
	fmt.Printf("My name is: %s; age: %d \n", p.Name, p.Age)
	time.Sleep(time.Second * 2)
	return nil
}

type SqliteDest struct {
	DB *sql.DB
}

func (s *SqliteDest) Save(ctx context.Context, p Person) error {
	query := `INSERT INTO person (name, age) VALUES (?, ?)`
	_, err := s.DB.ExecContext(ctx, query, p.Name, p.Age)
	return err
}

type Person struct {
	Name string `json:"person_name"`
	Age  int    `json:"person_age"`
}

func (p Person) HelloWorld() {
}

func main() {
	dsn := "file:data.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	createTableQuery := `CREATE TABLE IF NOT EXISTS person (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    age INTEGER);`
	if _, err := db.Exec(createTableQuery); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

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
	dest := SqliteDest{DB: db}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var people []Person
	err = json.Unmarshal([]byte(jsonString), &people)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	ch := make(chan Person)
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
					fmt.Printf("Inserted person: %s", p.Name)
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

	fmt.Println("All done. Check database")
}
