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

	for index, record := range records {
		if index == 1 {
			resources.HandleProcessAllKyc(w, r, record)
		}
	}
}
