package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kuo-52033/go-q/internal/model"
	"github.com/kuo-52033/go-q/internal/service"
	"github.com/redis/go-redis/v9"
)

type JobHandler interface {
	HandleJob(ctx context.Context, job *model.Job) error
}

type Manager struct {
	jobStore service.JobStore
	handlers map[string]JobHandler
	concurrency int
}

func NewManager(jobStore service.JobStore, concurrency int) *Manager {
	return &Manager{
		jobStore: jobStore,
		handlers: make(map[string]JobHandler),
		concurrency: concurrency,
	}
}

func (m *Manager) RegisterHandler(jobType string, handler JobHandler) {
	m.handlers[jobType] = handler
	log.Println("Registered handler for job type:", jobType)
}

func (m *Manager) StartWorker(ctx context.Context, queueName string) {
	log.Println("Starting worker for queue:", queueName)

	jobChan := make(chan string, m.concurrency)

	var workerWg sync.WaitGroup

	for i := 0; i < m.concurrency; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			m.processJobs(ctx, jobChan)
		}()
	}

	for{
		select {
			case <-ctx.Done():
				log.Println("Worker for queue", queueName, "stopped")
				close(jobChan)
				log.Println("Waiting for workers to finish...")
				workerWg.Wait()
				return
			default:
				jobId, err := m.jobStore.DequeueJobId(ctx, queueName, 5*time.Second)
				if err != nil {
					if err == redis.Nil {
						continue
					}
					log.Printf("Error dequeuing job: %v", err)
					// prevent busy loop
					time.Sleep(1 * time.Second)
					continue
				}
				if jobId != "" {
					jobChan <- jobId
				}
		}	
	}

}

func (m *Manager) processJobs(ctx context.Context, jobChan <-chan string) {
	for jobID := range jobChan {
		if err := m.processJob(ctx, jobID); err != nil {
			log.Printf("Error processing job %s: %v", jobID, err)
		}
	}
}

func (m *Manager) processJob(ctx context.Context, jobID string) error {
	job, err := m.jobStore.GetJobById(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	handler, ok := m.handlers[job.Type]
	if !ok {
		return fmt.Errorf("no handler found for job type: %s", job.Type)
	}

	if err := m.jobStore.UpdateJobStatus(ctx, jobID, model.StatusProcessing); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	log.Println("Processing job:", jobID)

	if err := handler.HandleJob(ctx, job); err != nil {
		//TODO: retry logic or mark as failed
		return fmt.Errorf("failed to handle job: %w", err)
	}

	if err := m.jobStore.UpdateJobStatus(ctx, jobID, model.StatusCompleted); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	log.Println("Job completed:", jobID)
	
	return nil
}
