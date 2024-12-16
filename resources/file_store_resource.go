package resources

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"trullio-kyc/exceptions"
	"trullio-kyc/utils"

	"github.com/xuri/excelize/v2"
)

var requestData struct {
	path     string
	FileName string `form:"file_name"`
}

type Record struct {
	ClientReferenceID        string    `json:"client_reference_id"`
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
}

func HandleCreateFile(w http.ResponseWriter, r *http.Request) (string, exceptions.ErrorResponse) {
	fmt.Println("Creating File")

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

	fmt.Println("Created file")
	return tempFile.Name(), exceptions.ErrorResponse{}
}

func StoreRecordsFromSpreadSheet(db *sql.DB, w http.ResponseWriter, r *http.Request, pathFile string) {
	// Mounting the Records
	records := ReadAndGetContentFile(pathFile)

	// Store Records From file
	StoreRecords(db, records)
}

func StoreRecords(db *sql.DB, records []Record) {

}

func ReadAndGetContentFile(pathFile string) []Record {
	fmt.Println("Reading File")
	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet_1")
	if err != nil {
		fmt.Println(err)
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

	return records
}
