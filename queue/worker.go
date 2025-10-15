package queue

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"trullio-kyc/config"
	"trullio-kyc/models"
	"trullio-kyc/resources"
)

type Job struct {
	ID     string
	Record models.Record
	W      http.ResponseWriter
	R      *http.Request
}

type Result struct {
	JobID string
	Error error
}

type WorkerPool struct {
	Workers     int
	JobQueue    chan Job
	ResultQueue chan Result
	quit        chan bool
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWorkerPool(workers int, jobQueueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		Workers:     workers,
		JobQueue:    make(chan Job, jobQueueSize),
		ResultQueue: make(chan Result, jobQueueSize),
		quit:        make(chan bool),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (wp *WorkerPool) Start() {
	config.AppLogger.Printf("Starting worker pool with %d workers", wp.Workers)

	for i := 0; i < wp.Workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	config.AppLogger.Println("Stopping worker pool...")
	wp.cancel()
	close(wp.quit)
	wp.wg.Wait()
	close(wp.JobQueue)
	close(wp.ResultQueue)
	config.AppLogger.Println("Worker pool stopped")
}

func (wp *WorkerPool) AddJob(job Job) {
	select {
	case wp.JobQueue <- job:
		config.AppLogger.Printf("Job %s added to queue", job.ID)
	case <-wp.ctx.Done():
		config.AppLogger.Printf("Worker pool stopped, cannot add job %s", job.ID)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	config.AppLogger.Printf("Worker %d started", id)

	for {
		select {
		case job := <-wp.JobQueue:
			config.AppLogger.Printf("Worker %d processing job %s for record ID %v", id, job.ID, job.Record.Id)

			start := time.Now()
			err := wp.processKYC(job)
			duration := time.Since(start)

			result := Result{
				JobID: job.ID,
				Error: err,
			}

			if err != nil {
				config.AppLogger.Printf("Worker %d failed to process job %s (took %v): %v", id, job.ID, duration, err)
			} else {
				config.AppLogger.Printf("Worker %d completed job %s for record ID %v (took %v)", id, job.ID, job.Record.Id, duration)
			}

			select {
			case wp.ResultQueue <- result:
			case <-wp.ctx.Done():
				return
			}

		case <-wp.quit:
			config.AppLogger.Printf("Worker %d stopping", id)
			return
		case <-wp.ctx.Done():
			config.AppLogger.Printf("Worker %d cancelled", id)
			return
		}
	}
}

func (wp *WorkerPool) processKYC(job Job) error {
	config.AppLogger.Printf("🚀 Worker starting isolated KYC processing for job %s, record ID %v", job.ID, job.Record.Id)

	// Add timeout for individual KYC processing
	ctx, cancel := context.WithTimeout(wp.ctx, 30*time.Second)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- resources.HandleProcessAllKyc(job.W, job.R, job.Record)
	}()

	select {
	case err := <-done:
		if err != nil {
			config.AppLogger.Printf("❌ Worker failed processing job %s: %v", job.ID, err)
		} else {
			config.AppLogger.Printf("✅ Worker completed job %s successfully", job.ID)
		}
		return err
	case <-ctx.Done():
		config.AppLogger.Printf("⏰ Worker timeout for job %s after 30 seconds", job.ID)
		return fmt.Errorf("KYC processing timeout for record ID %v", job.Record.Id)
	}
}

func (wp *WorkerPool) GetResults() <-chan Result {
	return wp.ResultQueue
}

// Rate limiter to control API calls to Trulioo
type RateLimiter struct {
	requests chan struct{}
	ticker   *time.Ticker
	quit     chan struct{}
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(chan struct{}, requestsPerSecond),
		ticker:   time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		quit:     make(chan struct{}),
	}

	// Fill initial bucket
	for i := 0; i < requestsPerSecond; i++ {
		rl.requests <- struct{}{}
	}

	go rl.refill()
	return rl
}

func (rl *RateLimiter) refill() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.requests <- struct{}{}:
			default:
				// Channel is full, skip
			}
		case <-rl.quit:
			rl.ticker.Stop()
			return
		}
	}
}

func (rl *RateLimiter) Wait() {
	<-rl.requests
}

func (rl *RateLimiter) Stop() {
	close(rl.quit)
}