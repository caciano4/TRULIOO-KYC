package controllers

import (
	"encoding/json"
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/models"
)

type ClientsResponse struct {
	Clients []models.Record `json:"clients"`
	Total   int             `json:"total"`
	Message string          `json:"message"`
}

func GetKYCClients(w http.ResponseWriter, r *http.Request) {
	config.AppLogger.Println("Fetching KYC clients list")

	// Connect to database
	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	// Query to get all clients with their KYC status
	query := `
		SELECT
			id, package_file_id, package_name, upload_by_id, client_reference_id, transfer_agent_responsible,
			type_of_transfer, email, user_id, first_name, middle_name, last_name, date_of_birth_day,
			personal_phone_number, street_address, city, postal, letter_state, letter_country,
			national_id, request, response, notes, match, complete_kyc, created_at, updated_at, deleted_at
		FROM public.document_records
		WHERE deleted_at IS NULL
		ORDER BY updated_at DESC, id DESC
		LIMIT 100
	`

	rows, err := db.Query(query)
	if err != nil {
		config.AppLogger.Printf("Error querying clients: %v", err)
		http.Error(w, "Error fetching clients", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []models.Record

	// Iterate through rows and append to clients slice
	for rows.Next() {
		var client models.Record
		err := rows.Scan(
			&client.Id,
			&client.PackageID,
			&client.PackageName,
			&client.UploadById,
			&client.ClientReferenceID,
			&client.TransferAgentResponsible,
			&client.TypeOfTransfer,
			&client.Email,
			&client.UserID,
			&client.FirstName,
			&client.MiddleName,
			&client.LastName,
			&client.DateOfBirthDay,
			&client.PersonalPhoneNumber,
			&client.StreetAddress,
			&client.City,
			&client.Postal,
			&client.LetterState,
			&client.LetterCountry,
			&client.NationalID,
			&client.Request,
			&client.Response,
			&client.Notes,
			&client.Match,
			&client.CompleteKYC,
			&client.CreatedAt,
			&client.UpdatedAt,
			&client.DeletedAt,
		)

		if err != nil {
			config.AppLogger.Printf("Error scanning client row: %v", err)
			continue
		}

		// Add status field based on processing state
		clients = append(clients, client)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		config.AppLogger.Printf("Error iterating client rows: %v", err)
		http.Error(w, "Error processing clients data", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := ClientsResponse{
		Clients: clients,
		Total:   len(clients),
		Message: "Clients fetched successfully",
	}

	config.AppLogger.Printf("Successfully fetched %d clients", len(clients))

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}