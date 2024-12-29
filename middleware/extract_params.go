package middleware

import (
	"context"
	"net/http"
	"strings"
	"trullio-kyc/config"
)

type ContextKey string

const ParamsKey ContextKey = "param"
const ParamIndex int = 2

func ExtractParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extractParam(next, w, r)
	})
}

func extractParam(next http.Handler, w http.ResponseWriter, r *http.Request) {
	pathParams := strings.Split(r.URL.Path, "/")

	if len(pathParams) > ParamIndex {
		// Extract param based on ParamIndex
		param := pathParams[ParamIndex]

		//Store on the context of request
		ctx := context.WithValue(r.Context(), ParamsKey, param)
		next.ServeHTTP(w, r.WithContext(ctx))
	} else {
		config.AppLogger.Println("Param not found.")
	}
}
