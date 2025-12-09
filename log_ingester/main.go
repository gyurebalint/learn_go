package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataDest interface {
	Save(p Person) error
}

type ConsoleDest struct{}

func (ConsoleDest) Save(p Person) error {
	if p.Name == "John Doe" {
		return errors.New("database connection lost")
	}
	fmt.Printf("My name is: %s; age: %d \n", p.Name, p.Age)
	time.Sleep(time.Second * 2)
	return nil
}

type Person struct {
	Name string `json:"person_name"`
	Age  int    `json:"person_age"`
}

func (p Person) HelloWorld() {
}

func main() {
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
	dest := DataDest(new(ConsoleDest))
	var people []Person
	err := json.Unmarshal([]byte(jsonString), &people)
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
			fmt.Println(err.Error())
		}
	}()

	//listener
	for range 3 {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			for p := range ch {
				if err := dest.Save(p); err != nil {
					errCh <- err
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
}
