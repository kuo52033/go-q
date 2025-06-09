package main

import (
	"fmt"
	"log"

	"github.com/kuo-52033/go-q/internal/db"
)

func main() {
	_, err := db.ConnectRedis("localhost:6379")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Redis connected successfully")

}
