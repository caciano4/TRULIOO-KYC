package config

import (
	"log"
	"os"
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
