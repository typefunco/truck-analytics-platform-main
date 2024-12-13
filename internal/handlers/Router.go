package handlers

import (
	"log"
	"net/http"
	"sync"

	october2023 "truck-analytics-platform/internal/handlers/2023/october"
	september2023 "truck-analytics-platform/internal/handlers/2023/september"
	october2024 "truck-analytics-platform/internal/handlers/2024/october"
	september2024 "truck-analytics-platform/internal/handlers/2024/september"
	"truck-analytics-platform/internal/handlers/utils"

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

		// 9 MONTH

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

		// Total market 9M 2023
		server.Handle("GET", "/9m2023tractors4x2total", september2023.Tractors4x2WithTotalMarket2023)
		server.Handle("GET", "/9m2023tractors6x4total", september2023.Tractors6x4WithTotalMarket2023)
		server.Handle("GET", "/9m2023dumpers6x4total", september2023.Dumpers6x4WithTotalMarket2023)
		server.Handle("GET", "/9m2023dumpers8x4total", september2023.Dumpers8x4WithTotalMarket2023)
		server.Handle("GET", "/9m2023ldttotal", september2023.NineMonth2023LDTTotal)
		server.Handle("GET", "/9m2023mdttotal", september2023.NineMonth2023MDTTotal)

		// -----------------------

		// Total market 9M 2024
		server.Handle("GET", "/9m2024tractors4x2total", september2024.Tractors4x2WithTotalMarket2024)
		server.Handle("GET", "/9m2024tractors6x4total", september2024.Tractors6x4WithTotalMarket2024)
		server.Handle("GET", "/9m2024dumpers6x4total", september2024.Dumpers6x4WithTotalMarket2024)
		server.Handle("GET", "/9m2024dumpers8x4total", september2024.Dumpers8x4WithTotalMarket2024)
		server.Handle("GET", "/9m2024ldttotal", september2024.NineMonth2024LDTTotal)
		server.Handle("GET", "/9m2024mdttotal", september2024.NineMonth2024MDTTotal)

		// -------------------------------------

		// MDT 2023 10M
		server.Handle("GET", "/10m2023mdt", october2023.TenMonth2023Mdt)
		server.Handle("GET", "/10m2023mdttotal", october2023.TenMonth2023MDTTotal)

		// LDT 2023 10M
		server.Handle("GET", "/10m2023ldt", october2023.TenMonth2023Ldt)
		server.Handle("GET", "/10m2023ldttotal", october2023.TenMonth2023LDTTotal)

		// MDT 2024 10M
		server.Handle("GET", "/10m2024mdt", october2024.TenMonth2024Mdt)
		server.Handle("GET", "/10m2024mdttotal", october2024.TenMonth2024MDTTotal)

		// LDT 2024 10M
		server.Handle("GET", "/10m2024ldt", october2024.TenMonth2024Ldt)
		server.Handle("GET", "/10m2024ldttotal", october2024.TenMonth2024LDTTotal)

		// HDT 2023 10M 4x2 Tractors
		server.Handle("GET", "/10m2023tractors4x2", october2023.TenMonth2023Tractors4x2)
		server.Handle("GET", "/10m2023tractors4x2total", october2023.TenTractors4x2WithTotalMarket2023)

		// HDT 2023 10M 6x4 Tractors
		server.Handle("GET", "/10m2023tractors6x4", october2023.TenMonth2023Tractors6x4)
		server.Handle("GET", "/10m2023tractors6x4total", october2023.TenTractors6x4WithTotalMarket2023)

		// HDT 2023 10M 6x4 Dumpers
		server.Handle("GET", "/10m2023dumpers6x4", october2023.TenMonth2023Dumpers6x4)
		server.Handle("GET", "/10m2023dumpers6x4total", october2023.TenDumpers6x4WithTotalMarket2023)

		// HDT 2023 10M 8x4 Dumpers
		server.Handle("GET", "/10m2023dumpers8x4", october2023.TenMonth2023Dumpers8x4)
		server.Handle("GET", "/10m2023dumpers8x4total", october2023.TenDumpers8x4WithTotalMarket2023)

		server.POST("/auth", AuthHandler)
		server.GET("/verify-token", VerifyTokenHandler)

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

// AuthHandler обрабатывает запросы на авторизацию
func AuthHandler(c *gin.Context) {
	var loginData struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	token, err := utils.CreateJWT(loginData.Login, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// VerifyTokenHandler проверяет валидность JWT токена
func VerifyTokenHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	err := utils.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}
