package resources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"trullio-kyc/config"
	"trullio-kyc/exceptions"
	"trullio-kyc/utils"

	"github.com/xuri/excelize/v2"
)

var requestData struct {
	path     string
	FileName string `form:"file_name"`
}

type Record struct {
	ClientReferenceID string `json:"client_reference_id"`
	// PackageName              string    `json:"package_name"`
	TransferAgentResponsible string    `json:"transfer_agent_responsible"`
	TypeOfTransfer           string    `json:"type_of_transfer"`
	Email                    string    `json:"email"`
	UserID                   string    `json:"user_id"`
	FirstName                string    `json:"first_name"`
	MiddleName               string    `json:"middle_name"`
	LastName                 string    `json:"last_name"`
	DateOfBirthDay           time.Time `json:"date_of_birth_day"`
	PersonalPhoneNumber      string    `json:"personal_phone_number"`
	StreetAddress            string    `json:"street_address"`
	City                     string    `json:"city"`
	Postal                   string    `json:"postal"`
	LetterState              string    `json:"letter_state"`
	LetterCountry            string    `json:"letter_country"`
	NationalID               string    `json:"national_id"`
	Request                  string    `json:"request"`
	Response                 string    `json:"response"`
	Notes                    string    `json:"notes"`
	Match                    string    `json:"match"`
	// packageID                string    `json:"package_id"`
}

func HandleCreateFile(w http.ResponseWriter, r *http.Request) (string, exceptions.ErrorResponse) {
	config.AppLogger.Println("Creating File")

	//Get Body
	requestData.FileName = r.FormValue("file_name")

	// Read and GET file
	file, _, err := r.FormFile("file")
	if err != nil {
		exceptions.NewErrorResponse("Failed to retrieve file", http.StatusBadRequest, err, w)
	}
	defer file.Close()

	// Create a Temporary file
	tempFile, err := os.CreateTemp(utils.GetProjectPath()+"/uploads", "upload-*.xlsx")
	if err != nil {
		exceptions.NewErrorResponse("Failed to create temporary file", http.StatusInternalServerError, err, w)
	}
	defer tempFile.Close()

	fileInBytes, err := io.ReadAll(file)
	if err != nil {
		exceptions.NewErrorResponse("Failed to get content from file", http.StatusInternalServerError, err, w)
	}
	//Write content file into a tempirary file
	if _, err := tempFile.Write(fileInBytes); err != nil {
		exceptions.NewErrorResponse("Failed to write to temporary file", http.StatusInternalServerError, err, w)
	}

	config.AppLogger.Println("Created file")
	return tempFile.Name(), exceptions.ErrorResponse{}
}

func StoreRecordsFromSpreadSheet(w http.ResponseWriter, r *http.Request, pathFile string) {
	// Mounting the Records
	records := ReadAndGetContentFile(pathFile)

	// Store Records From file
	StoreRecords(records, w, r)
}

func StoreRecords(records []Record, w http.ResponseWriter, r *http.Request) {
	var index int = 0
	packageFileId, err := utils.GenerateULIDWithDash()
	if err != nil {
		config.AppLogger.Printf("Failed to generate packageFileId %s", err.Error())
		return
	}
	config.AppLogger.Print(packageFileId)

	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//!TODO Create validation to check if already has that package_name
	// if so don't allow to upload again

	for _, record := range records {
		index++

		query := `
			INSERT INTO public.document_records 
			(package_file_id, package_name, upload_by_id, client_reference_id, transfer_agent_responsible, type_of_transfer, email, 
			user_id, first_name, middle_name, last_name, date_of_birth_day, personal_phone_number, 
			street_address, city, postal, letter_state, letter_country, national_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21);`
		_, err := db.Exec(
			query,
			packageFileId, requestData.FileName, 1, record.ClientReferenceID, record.TransferAgentResponsible,
			record.TypeOfTransfer, record.Email, record.UserID, record.FirstName,
			record.MiddleName, record.LastName, record.DateOfBirthDay,
			record.PersonalPhoneNumber, record.StreetAddress, record.City,
			record.Postal, record.LetterState, record.LetterCountry, record.NationalID,
			`NOW()`, `NOW()`,
		)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			config.AppLogger.Printf("Error inserting user: %v", err)
			return
		}

		//TODO: Check implement of second ifs
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			error := map[string]string{
				"error":       err.Error(),
				"failed_code": "Failed to insert User",
			}

			if jsonErr := json.NewEncoder(w).Encode(error); jsonErr != nil {
				config.AppLogger.Printf("Error encoding JSON: %v", jsonErr)
			}

			config.AppLogger.Printf("Error inserting User: %v", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": fmt.Sprintf("%d Records inserted successfully", index),
	}
	json.NewEncoder(w).Encode(response)
}

func ReadAndGetContentFile(pathFile string) []Record {
	config.AppLogger.Println("Reading File and mapping Records")

	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		config.AppLogger.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			config.AppLogger.Println(err)
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet_1")
	if err != nil {
		config.AppLogger.Println(err)
		return nil
	}

	// Map column headers to struct fields
	headerMap := make(map[string]int)
	var records []Record

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			// Map headers to their positions
			for colIndex, header := range row {
				headerMap[header] = colIndex
			}
			continue
		}

		// Parse each row into a Record struct
		record := Record{}

		// Set the UUID for ClientReferenceId
		record.ClientReferenceID = "manual-" + utils.GenerateUUID()

		for key, colIndex := range headerMap {
			if colIndex >= len(row) {
				continue
			}
			value := row[colIndex]

			switch key {
			case "TA Responsible":
				record.TransferAgentResponsible = value
			case "Type of Transfer":
				record.TypeOfTransfer = value
			case "email":
				record.Email = value
			case "USER ID":
				record.UserID = value
			case "First Name":
				record.FirstName = value
			case "Middle Name":
				record.MiddleName = value
			case "Last Name":
				record.LastName = value
			case "DOB (YYYY-MM-DD)":
				if dob, err := time.Parse("2006-01-02", value); err == nil {
					record.DateOfBirthDay = dob
				}
			case "Personal phone number":
				record.PersonalPhoneNumber = value
			case "Street Address":
				record.StreetAddress = value
			case "City":
				record.City = value
			case "Postal":
				record.Postal = value
			case "2 Letter State":
				record.LetterState = value
			case "2 Letter Country":
				record.LetterCountry = value
			case "National ID":
				record.NationalID = value
			case "Request":
				record.Request = value
			case "Response":
				record.Response = value
			case "Notes":
				record.Notes = value
			}
		}
		records = append(records, record)
	}

	config.AppLogger.Print("Finished to map Struct of Records")
	return records
}
