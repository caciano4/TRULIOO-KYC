package utils

import (
	"fmt"
	"os"
)

func GetProjectPath() string {
	// Get the current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %v\n", err)
		os.Exit(1)
	}
	return projectPath
}
