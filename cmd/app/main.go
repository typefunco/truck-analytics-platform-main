package main

import (
	"log/slog"
	"truck-analytics-platform/internal/db"
	"truck-analytics-platform/internal/handlers"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Connected to DB")

	handlers.InitRouter()
	slog.Info("Server started")
}
