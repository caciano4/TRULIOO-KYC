package controllers

import (
	"net/http"
	"trullio-kyc/exceptions"
	"trullio-kyc/resources"
	"trullio-kyc/validations"
)

func StoreFile(w http.ResponseWriter, r *http.Request) {
	//validate palyload
	validations.FileStoreValidate(r)

	// Store File
	pathFile, err := resources.HandleCreateFile(w, r)
	if err.Err != nil {
		exceptions.NewErrorResponse(err.Description, http.StatusInternalServerError, err.Err, w)
	}

	// resources.Handle
	resources.StoreRecordsFromSpreadSheet(w, r, pathFile)
}
