package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DebugRequest(feedBack *bool, w http.ResponseWriter, v ...interface{}) {
	// Default feedBack to false if nil
	if feedBack == nil {
		defaultValue := false
		feedBack = &defaultValue
	}

	var results []string
	for _, value := range v {
		switch v := value.(type) {
		case string:
			results = append(results, v) // Use the string as-is
		case []byte:
			results = append(results, string(v)) // Convert bytes to string
		default:
			jsonData, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				results = append(results, fmt.Sprintf("%v", value))
			} else {
				results = append(results, string(jsonData))
			}
		}
	}

	// Conditional feedback logic
	if *feedBack {
		fmt.Println("Feedback Enabled")
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Debug Data",
		"data":    results,
	})
}
