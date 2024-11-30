package handlers

import (
	"log"
	"net/http"
	"sync"

	september2023 "truck-analytics-platform/internal/handlers/2023/september"
	september2024 "truck-analytics-platform/internal/handlers/2024/september"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	var wg sync.WaitGroup

	// API-сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		server := gin.Default()
		server.Use(CORSMiddleware())

		// 2023
		// Tractors
		server.Handle("GET", "/9m2023tractors4x2", september2023.NineMonth2023Tractors4x2)
		server.Handle("GET", "/9m2023tractors6x4", september2023.NineMonth2023Tractors6x4)

		// Dumpers
		server.Handle("GET", "/9m2023dumpers6x4", september2023.NineMonth2023Dumpers6x4)
		server.Handle("GET", "/9m2023dumpers8x4", september2023.NineMonth2023Dumpers8x4)

		// LDT | MDT
		server.Handle("GET", "/9m2023ldt", september2023.NineMonth2023Ldt)
		server.Handle("GET", "/9m2023mdt", september2023.NineMonth2023Mdt)

		// -----------------------

		// 2024
		// Tractors
		server.Handle("GET", "/9m2024tractors4x2", september2024.NineMonth2024Tractors4x2)
		server.Handle("GET", "/9m2024tractors6x4", september2024.NineMonth2024Tractors6x4)

		// Dumpers
		server.Handle("GET", "/9m2024dumpers6x4", september2024.NineMonth2024Dumpers6x4)
		server.Handle("GET", "/9m2024dumpers8x4", september2024.NineMonth2024Dumpers8x4)

		// LDT | MDT
		server.Handle("GET", "/9m2024ldt", september2024.NineMonth2024Ldt)
		server.Handle("GET", "/9m2024mdt", september2024.NineMonth2024Mdt)

		// -----------------------

		// Total market 2023
		server.Handle("GET", "/9m2023tractors4x2total", september2023.Tractors4x2WithTotalMarket2023)
		server.Handle("GET", "/9m2023tractors6x4total", september2023.Tractors6x4WithTotalMarket2023)
		server.Handle("GET", "/9m2023dumpers6x4total", september2023.Dumpers6x4WithTotalMarket2023)
		server.Handle("GET", "/9m2023dumpers8x4total", september2023.Dumpers8x4WithTotalMarket2023)
		server.Handle("GET", "/9m2023ldttotal", september2023.NineMonth2023LDTTotal)
		server.Handle("GET", "/9m2023mdttotal", september2023.NineMonth2023MDTTotal)

		// -----------------------

		// Total market 2024
		server.Handle("GET", "/9m2024tractors4x2total", september2024.Tractors4x2WithTotalMarket2024)
		server.Handle("GET", "/9m2024tractors6x4total", september2024.Tractors6x4WithTotalMarket2024)
		server.Handle("GET", "/9m2024dumpers6x4total", september2024.Dumpers6x4WithTotalMarket2024)
		server.Handle("GET", "/9m2024dumpers8x4total", september2024.Dumpers8x4WithTotalMarket2024)
		server.Handle("GET", "/9m2024ldttotal", september2024.NineMonth2024LDTTotal)
		server.Handle("GET", "/9m2024mdttotal", september2024.NineMonth2024MDTTotal)

		log.Println("API server is running on port 8080...")
		if err := http.ListenAndServe(":8080", server); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Фронтенд-сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Frontend server is running on port 80...")
		if err := http.ListenAndServe(":80", http.FileServer(http.Dir("./frontend"))); err != nil {
			log.Fatalf("Failed to start frontend server: %v", err)
		}
	}()

	// Ожидание завершения работы серверов
	wg.Wait()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // завершает запрос на этапе OPTIONS
			return
		}

		c.Next()
	}
}
