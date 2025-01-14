package controllers

import (
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/resources"
)

func TruliooProcessingRequest(w http.ResponseWriter, r *http.Request) {
	//Catch all records to process KYC
	records, err := resources.HandleCatchKYCById(r)
	if err != nil {
		config.AppLogger.Fatal(err.Error())
	}

	// Iterate over records and enqueue each one for processing
	for _, record := range records {
		err := resources.HandleProcessAllKyc(w, r, record)
		if err != nil {
			config.AppLogger.Print("Failed to process record Id %v: %v", record.Id, err)
		} else {
			config.AppLogger.Print("Successfully queued record Id %v for processing", record.Id)
		}
	}
}
