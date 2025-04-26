package storage

import (
	"context"
	"effective-mobile-task/pkg/logger"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Psql *pgx.Conn
}

// Открываю коннект
func Open(url string) (*DB, error) {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		logger.Log.Errorf("Database connection error: %w", err)
		return nil, fmt.Errorf("Database connection error: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		logger.Log.Errorf("Error checking database connection: %v", err)
		return nil, fmt.Errorf("Error checking database connection: %w", err)
	}

	logger.Log.Println("Successfully connected to database")
	return &DB{Psql: conn}, nil
}

// Закрываю соединение
func (db *DB) Close() {
	if db.Psql != nil {
		db.Psql.Close(context.Background())
		logger.Log.Println("Closed database connection")
	}
}
