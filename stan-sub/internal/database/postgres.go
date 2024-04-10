package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/BatoBudaev/WB-L0/internal/model"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func InitDB(user, password, dbname string) (*DB, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8", user, password, dbname)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
    		data JSONB NOT NULL
		);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) InsertOrder(order model.Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO orders (data) VALUES ($1)`, jsonData)
	if err != nil {
		return err
	}

	return nil
}
