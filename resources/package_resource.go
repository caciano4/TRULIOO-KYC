package resources

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"trullio-kyc/config"
	"trullio-kyc/utils"
)

func HandleGetPackageList(w http.ResponseWriter, r *http.Request) {
	config.AppLogger.Println("Starting Search packages ...")
	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	query := `
		SELECT 
			COUNT(dr.id) AS total_records,
			MAX(package_name) AS package_name,
			SUM(CASE WHEN complete_kyc = true THEN 1 ELSE 0 END) AS completed,
			package_file_id AS package_id,
			CONCAT(u.first_name, ' ', u.last_name) AS full_name,
			MAX(transfer_agent_responsible) AS transfer_agent,
			MAX(type_of_transfer) AS type_of_transfer,
			MAX(dr.created_at) AS created,
			MAX(dr.updated_at) AS updated
		FROM document_records dr
		LEFT JOIN users u
			ON dr.upload_by_id = u.id
		WHERE dr.deleted_at IS NULL
		GROUP BY package_id, u.first_name, u.last_name
		ORDER BY created ASC
		`

	rows, err := db.Query(query)
	if err != nil {
		config.AppLogger.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to fetch packages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var packages []map[string]interface{}

	// Define nullable variables for scanning
	var totalRecords int
	var packageName, packageID, fullName, transferAgent, typeOfTransfer, completed, created, updated sql.NullString

	// Iterate over rows and build the response
	for rows.Next() {
		err := rows.Scan(
			&totalRecords,
			&packageName,
			&completed,
			&packageID,
			&fullName,
			&transferAgent,
			&typeOfTransfer,
			&created,
			&updated,
		)
		if err != nil {
			config.AppLogger.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing packages, Error message: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert nullable values to regular strings, using empty string for NULL values
		packageData := map[string]interface{}{
			"total_records":    totalRecords,
			"package_name":     getValue(packageName),
			"package_id":       getValue(packageID),
			"completed":        getValue(completed),
			"full_name":        getValue(fullName),
			"transfer_agent":   getValue(transferAgent),
			"type_of_transfer": getValue(typeOfTransfer),
			"created":          utils.FormatDate(getValue(created)),
			"updated":          utils.FormatDate(getValue(updated)),
		}

		packages = append(packages, packageData)
	}

	// Check for errors in rows iteration
	if err := rows.Err(); err != nil {
		config.AppLogger.Printf("Error iterating rows: %v", err)
		http.Error(w, "Error processing packages", http.StatusInternalServerError)
		return
	}

	// Build response
	response := map[string]interface{}{
		"message": "Success to fetch the packages",
		"data":    packages,
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "   ")
	encoder.Encode(response)
}

// Helper function to handle NULL values
func getValue(n sql.NullString) string {
	if n.Valid {
		return n.String
	}
	return ""
}
