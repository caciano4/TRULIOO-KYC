package validations

import (
	"fmt"
	"io"
	"net/http"
	"trullio-kyc/exceptions"

	"github.com/go-playground/validator/v10"
)

type Payload struct {
	File     []byte `form:"file" validate:"required,min=1"`
	FileName string `form:"file_name" validate:"required,min=1"`
}

func validatePayload(payload Payload) error {
	// Initialize the validator
	v := validator.New()

	// Validate the struct fields
	if err := v.Struct(payload); err != nil {
		return err
	}

	return nil
}

func FileStoreValidate(r *http.Request) {
	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		exceptions.NewErrorResponse("Failed to parse form data", http.StatusBadRequest, err, nil)
		return
	}

	// Extract file name and file
	fileName := r.FormValue("file_name")
	file, _, err := r.FormFile("file")
	if err != nil {
		exceptions.NewErrorResponse("Failed to retrieve file", http.StatusBadRequest, err, nil)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		exceptions.NewErrorResponse("Failed to read file content", http.StatusBadRequest, err, nil)
		return
	}

	// Create payload
	payload := Payload{
		File:     fileBytes,
		FileName: fileName,
	}

	// Validate the payload
	if err := validatePayload(payload); err != nil {
		exceptions.NewErrorResponse("Validation failed", http.StatusBadRequest, err, nil)
		return
	}

	fmt.Println("Payload is valid!")
}
