package controllers

import (
	"net/http"
	"trullio-kyc/resources"
)

func GetPackageList(w http.ResponseWriter, r *http.Request) {
	//
	resources.HandleGetPackageList(w, r)
}
