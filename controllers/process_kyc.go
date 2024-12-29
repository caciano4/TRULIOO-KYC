package controllers

import (
	"net/http"
	"trullio-kyc/config"
)

func InitTrulioo(w http.ResponseWriter, r *http.Request) {
	
	value := r.Context().Value("param").(string)
	config.AppLogger.Print("Chegou aqui!" + value)
}
