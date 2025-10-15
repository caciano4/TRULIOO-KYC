package controllers

import (
	"encoding/json"
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/queue"
	"trullio-kyc/resources"
)

type KYCProcessResponse struct {
	Message    string   `json:"message"`
	TotalJobs  int      `json:"total_jobs"`
	JobIDs     []string `json:"job_ids"`
	QueueStats map[string]interface{} `json:"queue_stats"`
}

func TruliooProcessingRequest(w http.ResponseWriter, r *http.Request) {
	// Initialize queue if not already done
	queue.InitializeQueue()
	qm := queue.GetQueueManager()

	// Catch all records to process KYC
	records, err := resources.HandleCatchKYCById(r)
	if err != nil {
		config.AppLogger.Printf("Error fetching KYC records: %v", err)
		http.Error(w, "Failed to fetch KYC records", http.StatusInternalServerError)
		return
	}

	if len(records) == 0 {
		response := KYCProcessResponse{
			Message:    "No records found for processing",
			TotalJobs:  0,
			JobIDs:     []string{},
			QueueStats: qm.GetQueueStats(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var jobIDs []string
	successCount := 0

	// Iterate over records and enqueue each one for processing
	for _, record := range records {
		jobID, err := qm.AddKYCJob(record, w, r)
		if err != nil {
			config.AppLogger.Printf("Failed to queue record Id %v: %v", record.Id, err)
		} else {
			jobIDs = append(jobIDs, jobID)
			successCount++
			config.AppLogger.Printf("Successfully queued record Id %v with job ID %s", record.Id, jobID)
		}
	}

	// Return response with job information
	response := KYCProcessResponse{
		Message:    "KYC processing jobs queued successfully",
		TotalJobs:  successCount,
		JobIDs:     jobIDs,
		QueueStats: qm.GetQueueStats(),
	}

	config.AppLogger.Printf("Queued %d out of %d records for KYC processing", successCount, len(records))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 - Accepted for async processing
	json.NewEncoder(w).Encode(response)
}
