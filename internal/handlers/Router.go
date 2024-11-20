package handlers

import (
	"net/http"

	september2023 "truck-analytics-platform/internal/handlers/2023/september"
	september2024 "truck-analytics-platform/internal/handlers/2024/september"

	"github.com/gin-gonic/gin"
)

func InitRouter() {

	server := gin.Default()
	server.Use(CORSMiddleware())

	// 2023

	// Tractors
	server.Handle("GET", "/9m2023tractors4x2", september2023.NineMonth2023Tractors4x2)
	server.Handle("GET", "/9m2023tractors6x4", september2023.NineMonth2023Tractors6x4)

	// Dumpers
	server.Handle("GET", "/9m2023dumpers6x4", september2023.NineMonth2023Dumpers6x4)
	server.Handle("GET", "/9m2023dumpers8x4", september2023.NineMonth2023Dumpers8x4)

	// -----------------------

	// 2024

	// Tractors
	server.Handle("GET", "/9m2024tractors4x2", september2024.NineMonth2024Tractors4x2)
	server.Handle("GET", "/9m2024tractors6x4", september2024.NineMonth2024Tractors6x4)

	// Dumpers
	server.Handle("GET", "/9m2024dumpers6x4", september2024.NineMonth2024Dumpers6x4)
	server.Handle("GET", "/9m2024dumpers8x4", september2024.NineMonth2024Dumpers8x4)

	// Total market 2023
	server.Handle("GET", "/9m2023tractors4x2total", september2023.Tractors4x2WithTotalMarket2023)
	server.Handle("GET", "/9m2023tractors6x4total", september2023.Tractors6x4WithTotalMarket2023)
	server.Handle("GET", "/9m2023dumpers6x4total", september2023.Dumpers6x4WithTotalMarket2023)
	server.Handle("GET", "/9m2023dumpers8x4total", september2023.Dumpers8x4WithTotalMarket2023)

	// Total market 2024
	server.Handle("GET", "/9m2024tractors4x2total", september2024.Tractors4x2WithTotalMarket2024)
	server.Handle("GET", "/9m2024tractors6x4total", september2024.Tractors6x4WithTotalMarket2024)
	server.Handle("GET", "/9m2024dumpers6x4total", september2024.Dumpers6x4WithTotalMarket2024)
	server.Handle("GET", "/9m2024dumpers8x4total", september2024.Dumpers8x4WithTotalMarket2024)

	http.ListenAndServe(":8080", server)

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
