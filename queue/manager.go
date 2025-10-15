package queue

import (
	"fmt"
	"net/http"
	"sync"
	"trullio-kyc/config"
	"trullio-kyc/models"
)

var (
	GlobalWorkerPool *WorkerPool
	GlobalRateLimiter *RateLimiter
	once             sync.Once
)

type QueueManager struct {
	workerPool    *WorkerPool
	rateLimiter   *RateLimiter
	jobCounter    int64
	mu            sync.Mutex
}

var queueManager *QueueManager

func InitializeQueue() {
	once.Do(func() {
		workers := config.GetEnvInt("WORKER_POOL_SIZE", 5)
		queueSize := config.GetEnvInt("JOB_QUEUE_SIZE", 100)
		rateLimit := config.GetEnvInt("TRULIOO_RATE_LIMIT", 10) // requests per second

		config.AppLogger.Printf("Initializing queue manager with %d workers, queue size %d, rate limit %d/s",
			workers, queueSize, rateLimit)

		GlobalWorkerPool = NewWorkerPool(workers, queueSize)
		GlobalRateLimiter = NewRateLimiter(rateLimit)

		queueManager = &QueueManager{
			workerPool:  GlobalWorkerPool,
			rateLimiter: GlobalRateLimiter,
		}

		// Start the worker pool
		GlobalWorkerPool.Start()

		// Start result processor
		go queueManager.processResults()

		config.AppLogger.Println("Queue manager initialized successfully")
	})
}

func GetQueueManager() *QueueManager {
	if queueManager == nil {
		InitializeQueue()
	}
	return queueManager
}

func (qm *QueueManager) AddKYCJob(record models.Record, w http.ResponseWriter, r *http.Request) (string, error) {
	qm.mu.Lock()
	qm.jobCounter++
	jobID := fmt.Sprintf("kyc-%d-%d", record.Id, qm.jobCounter)
	qm.mu.Unlock()

	job := Job{
		ID:     jobID,
		Record: record,
		W:      w,
		R:      r,
	}

	// Apply rate limiting before adding to queue
	qm.rateLimiter.Wait()

	qm.workerPool.AddJob(job)
	return jobID, nil
}

func (qm *QueueManager) processResults() {
	config.AppLogger.Println("Starting result processor")

	for result := range qm.workerPool.GetResults() {
		if result.Error != nil {
			config.AppLogger.Printf("Job %s completed with error: %v", result.JobID, result.Error)
			// Here you could implement retry logic or error handling
		} else {
			config.AppLogger.Printf("Job %s completed successfully", result.JobID)
		}
	}

	config.AppLogger.Println("Result processor stopped")
}

func (qm *QueueManager) GetQueueStats() map[string]interface{} {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	return map[string]interface{}{
		"workers":         qm.workerPool.Workers,
		"jobs_in_queue":   len(qm.workerPool.JobQueue),
		"results_pending": len(qm.workerPool.ResultQueue),
		"total_jobs":      qm.jobCounter,
	}
}

func (qm *QueueManager) Shutdown() {
	config.AppLogger.Println("Shutting down queue manager...")

	if qm.workerPool != nil {
		qm.workerPool.Stop()
	}

	if qm.rateLimiter != nil {
		qm.rateLimiter.Stop()
	}

	config.AppLogger.Println("Queue manager shutdown complete")
}

// Graceful shutdown handler
func SetupGracefulShutdown() {
	// This can be called from main() to setup proper shutdown handling
	// when the application receives termination signals
}