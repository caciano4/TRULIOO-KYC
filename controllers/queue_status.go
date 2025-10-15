package controllers

import (
	"encoding/json"
	"net/http"
	"trullio-kyc/queue"
)

type QueueStatusResponse struct {
	Status     string                 `json:"status"`
	QueueStats map[string]interface{} `json:"queue_stats"`
}

func GetQueueStatus(w http.ResponseWriter, r *http.Request) {
	qm := queue.GetQueueManager()

	response := QueueStatusResponse{
		Status:     "active",
		QueueStats: qm.GetQueueStats(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}