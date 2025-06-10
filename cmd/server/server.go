package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kuo-52033/go-q/internal/api/routes"
	"github.com/kuo-52033/go-q/internal/db"
	"github.com/kuo-52033/go-q/internal/api/middleware/common"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_URL")
	serverPort := os.Getenv("SERVER_PORT")

	rdb, err := db.NewRedisClient(redisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer rdb.Close()

	log.Println("Redis connected successfully")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.Use(common.ErrorHandler())

	api := router.Group("/api/v1")

	routes.SetupJobModule(api, rdb)

	log.Printf("Server is running on port %s", serverPort)
	router.Run(fmt.Sprintf(":%s", serverPort))
}
