package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/routes"
	"github.com/kuo-52033/go-q/internal/db"
)

func main() {
	rdb, err := db.NewRedisClient("localhost:6379")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer rdb.Close()

	fmt.Println("Redis connected successfully")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")

	routes.SetupJobModule(api, rdb)

	log.Println("Server is running on port 8080")
	router.Run(":8080")
}
