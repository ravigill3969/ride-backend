package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlconn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL!")

	return DB, nil
}
