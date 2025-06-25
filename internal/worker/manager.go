package worker

import (
	"context"
	"errors"
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
		jobId, err := m.jobStore.DequeueJobId(ctx, queueName, 5*time.Second)
		if err != nil {
			log.Println("get jobId error:", err)
			if errors.Is(err, redis.Nil) {

				select {
					case <-ctx.Done():
						goto cleanUp
					default:
						continue
				}
			}

			log.Printf("Error dequeuing job from '%s': %v. Retrying...", queueName, err)

			select {
				case <-time.After(5 * time.Second):
					continue
				case <-ctx.Done():
					goto cleanUp
			}
		}

		if jobId == "" {
			continue
		}

		select {
			case jobChan <- jobId:
				log.Println("Dequeued job:", jobId)
			case <-ctx.Done():
				log.Println("Context done, re-enqueueing job:", jobId)
				if err := m.jobStore.ReEnqueueJobId(ctx, queueName, jobId); err != nil {
					log.Printf("Failed to re-enqueue job %s: %v", jobId, err)
				}
				goto cleanUp
		}
	}

	cleanUp:
		close(jobChan)
		log.Println("Waiting for workers to finish...")
		workerWg.Wait()
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
