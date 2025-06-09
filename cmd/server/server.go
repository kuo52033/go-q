package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/handler/job"
	"github.com/kuo-52033/go-q/internal/api/routes"
	"github.com/kuo-52033/go-q/internal/db"
	"github.com/kuo-52033/go-q/internal/service"
)

func main() {
	rdb, err := db.NewRedisClient("localhost:6379")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer rdb.Close()

	fmt.Println("Redis connected successfully")

	jobService := service.NewJobService(rdb)
	jobHandler := job.NewHandler(jobService)

	router := gin.Default()

	routes.SetupRoutes(router, &routes.RouteConfig{
		JobHandler: jobHandler,
	})

	router.Run(":8080")
}
