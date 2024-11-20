package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	//connectionURL := "postgresql://postgres:postgres@localhost:5432/truck-analytics"
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	connect, err := pgx.Connect(context.Background(), connectionURL)
	if err != nil {
		slog.Error("Can't connect to DB")
		return nil, err
	}

	err = connect.Ping(context.Background())
	if err != nil {
		slog.Error("Can't Ping DB")
		return nil, err
	}

	slog.Info("Connected to DB")

	return connect, nil
}
