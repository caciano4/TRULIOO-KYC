package resources

import (
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
		ORDER BY created ASC`

	rows, err := db.Query(query)
	if err != nil {
		config.AppLogger.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to fetch packages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var packages []map[string]interface{}

	// Iterate over rows and build the response
	for rows.Next() {
		var totalRecords int
		var packageName, packageID, fullName, transferAgent, typeOfTransfer, created, updated string

		err := rows.Scan(
			&totalRecords,
			&packageName,
			&packageID,
			&fullName,
			&transferAgent,
			&typeOfTransfer,
			&created,
			&updated,
		)
		if err != nil {
			config.AppLogger.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing packages", http.StatusInternalServerError)
			return
		}

		// Add row to response
		packages = append(packages, map[string]interface{}{
			"total_records":    totalRecords,
			"package_name":     packageName,
			"package_id":       packageID,
			"full_name":        fullName,
			"transfer_agent":   transferAgent,
			"type_of_transfer": typeOfTransfer,
			"created":          utils.FormatDate(created),
			"updated":          utils.FormatDate(updated),
		})
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
