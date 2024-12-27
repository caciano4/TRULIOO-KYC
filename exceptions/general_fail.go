package exceptions

import (
	"encoding/json"
	"net/http"
	"trullio-kyc/config"
)

// ErrorResponse define o formato da resposta de erro
type ErrorResponse struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Error       string `json:"error"`
	Err         error
}

// NewErrorResponse cria uma nova estrutura de erro genérica
func NewErrorResponse(description string, statusCode int, err error, w http.ResponseWriter) ErrorResponse {
	errResponse := ErrorResponse{
		StatusCode:  statusCode,
		Description: description,
		Error:       err.Error(),
		Err:         err,
	}
	SendErrorResponse(w, errResponse)
	return errResponse
}

// SendErrorResponse envia a resposta de erro como JSON no HTTP
func SendErrorResponse(w http.ResponseWriter, errResponse ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResponse.StatusCode)
	json.NewEncoder(w).Encode(errResponse)
	config.AppLogger.Printf("Error: %s, Status: %d", errResponse.Err.Error(), errResponse.StatusCode)
}
