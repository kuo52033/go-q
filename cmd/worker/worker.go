package main

import (
	"context"
	"log"
	"os"
	"sync"
	"os/signal"
	"syscall"
	"time"
	"github.com/joho/godotenv"
	"github.com/kuo-52033/go-q/internal/db"
	"github.com/kuo-52033/go-q/internal/store"
	"github.com/kuo-52033/go-q/internal/worker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_URL")

	rdb, err := db.NewRedisClient(redisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	
	defer rdb.Close()

	log.Println("Redis connected successfully")

	jobStore := store.NewRedisJobStore(rdb)

	manager := worker.NewManager(jobStore, 3)

	manager.RegisterHandler("process_image", &worker.ImageHandler{})
	manager.RegisterHandler("send_email", &worker.EmailHandler{})
	manager.RegisterHandler("generate_report", &worker.ReportHandler{})

	queue := []string{"process_image", "send_email", "generate_report"}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for _, q := range queue {
		wg.Add(1)
		go func(q string) {
			defer wg.Done()
			manager.StartWorker(ctx, q)
		}(q)
	}

	log.Printf("workers started, listening on queues: %v", queue)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	cancel()

	log.Println("shutting down...")

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
		case <-done:
			log.Println("workers stopped")
		case <-time.After(10 * time.Second):
			log.Println("workers stopped with timeout")
	}
}
