package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"trullio-kyc/utils"
)

// Global logger variable
var AppLogger *log.Logger

func init() {
	path := utils.GetProjectPath() + "/log/" + time.Now().Format("2006-01-02") + ".log"

	// Initialize the global logger
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}

	// Create a new logger
	AppLogger = log.New(file, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogResponseTrulio(step int, user_name string, response interface{}, type_folder string) {
	var step_name string

	switch step {
	case 1:
		step_name = "step_1"
	case 2:
		step_name = "step_2"
	case 3:
		step_name = "step_3"
	case 4:
		step_name = "step_4"
	default:
		step_name = "step_not_found"
	}

	// Construct file path
	dirPath := fmt.Sprintf("%s/responses/%s/%s", utils.GetProjectPath(), strings.ToLower(user_name), type_folder)
	filePath := fmt.Sprintf("%s/%s_%s.json", dirPath, time.Now().Format("2006-01-02"), step_name)

	// Ensure the directory exists
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		AppLogger.Printf("Failed to create directory: %v", err)
		return
	}

	// Open or create the file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		AppLogger.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	// Marshal the response into JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		AppLogger.Printf("Failed to marshal response: %v", err)
		return
	}

	// Write JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		AppLogger.Printf("Failed to write to file: %v", err)
		return
	}

	AppLogger.Printf("Response logged to %s", filePath)
}
