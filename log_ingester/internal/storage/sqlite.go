package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log_ingester/internal/models"
)

type DataDest interface {
	Save(ctx context.Context, p models.Person) error
}

type SqliteDest struct {
	DB *sql.DB
}

func NewSqliteDest(dbUrl string) (*SqliteDest, error) {
	dsn := "file:data.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS person (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    age INTEGER);`
	if _, err := db.Exec(createTableQuery); err != nil {
		fmt.Println("Error creating table:", err)
		return nil, err
	}

	return &SqliteDest{db}, nil
}

func (s *SqliteDest) Save(ctx context.Context, p models.Person) error {
	query := `INSERT INTO person (name, age) VALUES (?, ?)`
	_, err := s.DB.ExecContext(ctx, query, p.Name, p.Age)
	return err
}

func (s *SqliteDest) Close() error {
	return s.DB.Close()
}
