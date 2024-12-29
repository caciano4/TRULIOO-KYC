package controllers

import (
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/middleware"
)

func InitTrulioo(w http.ResponseWriter, r *http.Request) {

	param := r.Context().Value(middleware.ParamsKey)

	config.AppLogger.Print(param)
}
