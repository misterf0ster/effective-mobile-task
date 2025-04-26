package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Psql *pgx.Conn
}

// Открываю коннект
func Open(url string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("Database connection error: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("Error checking database connection: %w", err)
	}

	log.Println("Successfully connected to database")
	return &DB{Psql: conn}, nil
}

// Закрываю соединение
func (db *DB) Close() {
	if db.Psql != nil {
		db.Psql.Close(context.Background())
		log.Println("Closed database connection")
	}
}
