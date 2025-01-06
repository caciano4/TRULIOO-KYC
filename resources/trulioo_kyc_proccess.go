package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"trullio-kyc/config"
	"trullio-kyc/middleware"
	"trullio-kyc/models"
	"trullio-kyc/utils"
)

type Req struct {
	FlowId string
	URL    string
	Body   map[string]interface{}
}

type TField struct {
	ID     string
	Name   string
	TValue string
}

var xHfSession string
var bearerToken string

type Fields []TField

func HandleCatchKYCById(r *http.Request) ([]models.Record, error) {
	// Track Log
	config.AppLogger.Print("Get KYC RECORDS")

	// Start variables
	param := r.Context().Value(middleware.ParamsKey).(string)
	var records []models.Record

	// Connect and Close database
	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	// Preparing query to fetch KYC by package_file_id
	query := `SELECT 
		id, package_file_id, package_name, upload_by_id, client_reference_id, transfer_agent_responsible, 
		type_of_transfer, email, user_id, first_name, middle_name, last_name, date_of_birth_day, 
		personal_phone_number, street_address, city, postal, letter_state, letter_country, 
		national_id, request, response, notes, match, complete_kyc, created_at, updated_at, deleted_at
	FROM public.document_records
	WHERE package_file_id = $1
	AND deleted_at IS NULL`

	// Use parameterized query to avoid SQL injection
	rows, err := db.Query(query, param)
	if err != nil {
		config.AppLogger.Println(err.Error())
		return records, nil
	}
	defer rows.Close()

	// Iterate through rows and append to records slice
	for rows.Next() {
		var record models.Record
		err := rows.Scan(
			&record.Id,
			&record.PackageID,
			&record.PackageName,
			&record.ClientReferenceID,
			&record.UploadById,
			&record.TransferAgentResponsible,
			&record.TypeOfTransfer,
			&record.Email,
			&record.UserID,
			&record.FirstName,
			&record.MiddleName,
			&record.LastName,
			&record.DateOfBirthDay,
			&record.PersonalPhoneNumber,
			&record.StreetAddress,
			&record.City,
			&record.Postal,
			&record.LetterState,
			&record.LetterCountry,
			&record.NationalID,
			&record.Request,
			&record.Response,
			&record.Notes,
			&record.Match,
			&record.CompleteKYC, // New field
			&record.CreatedAt,   // New field
			&record.UpdatedAt,   // New field
			&record.DeletedAt,   // New field
		)

		if err != nil {
			config.AppLogger.Printf("Error scanning row: %v", err)
			return records, nil
		}

		// Append to records slice
		records = append(records, record)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		config.AppLogger.Printf("Error iterating rows: %v", err)
		return records, nil
	}

	// Return records
	return records, nil
}

func HandleProcessAllKyc(w http.ResponseWriter, r *http.Request, record models.Record) {
	// Init and catch field IDS
	fields := truliooInit(w, record)

	// Send body with IDS, and store the request sent
	truliooBodySubmit(fields, record)

	// Retrieve Bearer Token
	truliooGenerateBearerToken(record)

	// MatchApi Trulioo
	truliooDetailsFromClient(w, r, record)
}

// step 4
func truliooDetailsFromClient(w http.ResponseWriter, r *http.Request, record models.Record) {
	config.AppLogger.Print("TRULIOO DETAILS FROM CLIENT: STEP 4")
	var clientDetails models.ClientDetailsResponse
	var request Req
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)
	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	request.URL = fmt.Sprintf("https://api.workflow.prod.trulioo.com/export/test/v2/query/client/%s?includeFullServiceDetails=true", xHfSession)

	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	if err := json.Unmarshal(body, &clientDetails); err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	config.LogResponseTrulio(4, userName, clientDetails, "response")

	//! TODO CATCH MATCH DATA
	query := `UPDATE document_records
				SET match = $1
			 WHERE id = $2`

	json, err := json.Marshal(clientDetails.FlowData)
	if err != nil {
		config.AppLogger.Print(fmt.Sprintf("Error serializing FlowData: %v", err))
	}

	result, err := db.Exec(query, json, *&record.Id)
	if err != nil {
		config.AppLogger.Print(err.Error())
	}

	config.AppLogger.Print(result)
}

// step 3
func truliooGenerateBearerToken(record models.Record) {
	config.AppLogger.Print("TRULIOO GENERATE BEARER TOKEN: STEP 3")
	var request Req
	var bearerTokenResponse models.BearerTokenReponse
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)

	//Using that type of body (x-www-form-urlencoded) because it's required as an oauth2 api
	request.URL = "https://auth-api.trulioo.com/connect/token"
	payload := strings.NewReader(
		fmt.Sprintf(
			"client_id=%s&client_secret=%s&grant_type=client_credentials",
			config.GetEnv("CLIENT_ID", ""),
			config.GetEnv("CLIENT_SECRET", ""),
		),
	)

	req, err := http.NewRequest("POST", request.URL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", "incap_ses_672_2454916=B5Qffro331hxKA5UpmxTCbpeeWcAAAAAy+cvmxfB6K02m855TXWcSQ==; visid_incap_2454916=y2SBkmdITHSDGY2dFArmi//3Y2cAAAAAQUIPAAAAAABxADON0lB0WMtw2kKg4f/O")

	client := &http.Client{}
	res, _ := client.Do(req)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	if err := json.Unmarshal(body, &bearerTokenResponse); err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	config.LogResponseTrulio(3, userName, bearerTokenResponse, "response")

	bearerToken = bearerTokenResponse.AccessToken
}

// step 2
func truliooBodySubmit(fields Fields, record models.Record) {
	config.AppLogger.Print("TRULIOO SUBMIT: STEP 2")
	var request Req
	var truliooBodySubmitResponse models.DirectSubmitResponse
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)
	if request.Body == nil {
		request.Body = make(map[string]interface{})
	}

	request.FlowId = config.GetEnv("FLOW_ID", "")
	request.URL = fmt.Sprintf("https://api.workflow.prod.trulioo.com/interpreter-v2/test/submit/%s", request.FlowId)

	for _, field := range fields {
		if field.ID != "" && field.TValue != "" {
			request.Body[field.ID] = field.TValue
		}
	}

	// Convert body map json
	bodyJson, err := json.Marshal(request.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	config.LogResponseTrulio(2, userName, request.Body, "request")

	req, err := http.NewRequest("POST", request.URL, bytes.NewBuffer(bodyJson))
	if err != nil {
		config.AppLogger.Printf(err.Error())
	}

	// Setting Header content type to json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hf-retry-on-pending", "true")

	// Submiting KYC request step 2
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		config.AppLogger.Print(err.Error())
	}
	defer res.Body.Close()

	//GET RESPONSE FROM STEP 2o
	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return
	}

	// Decoding json HTTP
	if err := json.Unmarshal(body, &truliooBodySubmitResponse); err != nil {
		config.AppLogger.Print(err.Error())
	}

	config.LogResponseTrulio(2, userName, truliooBodySubmitResponse, "response")

	//! TODO Add resposne return from trulioo in DB
	//! TODO Update ROW with response text
	// truliooBodySubmitResponse.Text

	// GETTING XHFSESSION
	xHfSession = res.Header.Get("x-hf-session")
}

// ! Step 1
func truliooInit(w http.ResponseWriter, record models.Record) Fields {
	config.AppLogger.Print("INIT TRULIOO REQUEST: STEP 1")

	// Variables
	var request Req
	var initTrulioo models.InitTrulioo
	var fields Fields
	var tField TField
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)

	// Preparing Req struct
	request.FlowId = config.GetEnv("FLOW_ID", "")
	request.URL = fmt.Sprintf("https://api.workflow.prod.trulioo.com/interpreter-v2/test/flow/%s", request.FlowId)

	// Instance new Request
	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		config.AppLogger.Print(err.Error())
	}

	// Adding the header Params
	req.Header.Add("Cookie", "visid_incap_2454916=y2SBkmdITHSDGY2dFArmi//3Y2cAAAAAQUIPAAAAAABxADON0lB0WMtw2kKg4f/O")

	// Start request and defer cloing the body
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		config.AppLogger.Print(err.Error())
	}
	defer res.Body.Close()

	// Reading the body response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
	}

	// Decoding Json response (se aplicável)
	if err := json.Unmarshal(body, &initTrulioo); err != nil {
		http.Error(w, "Error to decode json", http.StatusInternalServerError)
		config.AppLogger.Print(err.Error())
		return Fields{}
	}

	//Log The response
	config.LogResponseTrulio(1, userName, initTrulioo, "response")

	// Catch the IDS
	for _, element := range initTrulioo.Elements {
		if element.Role == "external_customer_id" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.ClientReferenceID, "")
		}

		if element.Role == "address_country" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.LetterCountry, "")
		}

		if element.Role == "first_name" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.FirstName, "")
		}

		if element.Role == "last_name" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.LastName, "")
		}

		if record.MiddleName != nil {
			if element.NormalizedName == "MiddleName" {
				tField.ID = element.ID
				tField.Name = element.NormalizedName
				tField.TValue = utils.GetStringValue(record.MiddleName, "")
			}
		}

		if element.Role == "dob" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = record.DateOfBirthDay.Format("2006-01-02")
		}

		if element.Role == "address_1" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.StreetAddress, "")
		}

		if element.Role == "address_city" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.City, "")
		}

		if element.NormalizedName == "Suburb" {
			mandatoryCountryToSuburb := map[string]struct{}{
				"AU": {}, "CA": {}, "DO": {}, "HK": {}, "KR": {}, "NO": {}, "PH": {}, "US": {}, "VE": {},
			}

			if _, exists := mandatoryCountryToSuburb[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.NormalizedName
				tField.TValue = utils.GetStringValue(record.Suburb, "")
			}
		}

		if element.Role == "address_state" {
			mandatoryCountryToState := map[string]struct{}{
				"AR": {}, "AU": {}, "BR": {}, "CA": {}, "CL": {}, "CO": {}, "CR": {}, "DO": {}, "EC": {}, "GR": {},
				"HK": {}, "IN": {}, "IT": {}, "KR": {}, "MY": {}, "MX": {},
			}

			if _, exists := mandatoryCountryToState[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.Role
				tField.TValue = utils.GetStringValue(record.LetterState, "")
			}
		}

		if element.Role == "address_zip" {
			tField.ID = element.ID
			tField.Name = element.Role
			tField.TValue = utils.GetStringValue(record.Postal, "")
		}

		if element.NormalizedName == "DriverLicenceNumber" {
			mandatoryCountryToDriverLicense := map[string]struct{}{
				"IN": {}, "NZ": {},
			}

			if _, exists := mandatoryCountryToDriverLicense[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.NormalizedName
				tField.TValue = utils.GetStringValue(record.DriverLicence, "")
			}
		}

		if element.NormalizedName == "DriverLicenceVersionNumber" {
			mandatoryCountryToDriverLicenseVersionNumber := map[string]struct{}{
				"NZ": {},
			}

			if _, exists := mandatoryCountryToDriverLicenseVersionNumber[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.NormalizedName
				tField.TValue = utils.GetStringValue(record.DriverLicenceVersion, "")
			}
		}

		if element.Role == "VoterID" {
			mandatoryCountryToVoterID := map[string]struct{}{
				"IN": {}, "GH": {}, "NG": {},
			}

			if _, exists := mandatoryCountryToVoterID[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.Role
				tField.TValue = utils.GetStringValue(record.VoterID, "")
			}
		}

		if element.Role == "social_service_number" {
			mandatoryCountryToSocialNumber := map[string]struct{}{
				"CA": {}, "US": {}, "PH": {},
			}

			if _, exists := mandatoryCountryToSocialNumber[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.Role
				tField.TValue = *record.NationalID
			}
		}

		if element.NormalizedName == "PassportNumber" {
			mandatoryCountryToPassportNumber := map[string]struct{}{
				"KE": {}, "GH": {},
			}

			if _, exists := mandatoryCountryToPassportNumber[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.NormalizedName
				tField.TValue = utils.GetStringValue(record.Passport, "")
			}
		}

		if element.Role == "national_id_nr" {
			mandatoryCountryToNationalId := map[string]struct{}{
				"AR": {}, "BH": {}, "BD": {}, "BR": {}, "CN": {},
				"CO": {}, "CR": {}, "DK": {}, "DO": {}, "EC": {},
				"EG": {}, "SV": {}, "GR": {}, "IS": {}, "IN": {},
				"IT": {}, "JO": {}, "KE": {}, "KW": {}, "LV": {},
				"LB": {}, "MT": {}, "MY": {}, "MX": {}, "NG": {},
				"OM": {}, "PE": {}, "QA": {}, "RO": {}, "SA": {},
				"ZA": {}, "ES": {}, "SE": {}, "TH": {}, "UA": {},
				"UY": {}, "VE": {},
			}

			if _, exists := mandatoryCountryToNationalId[*record.LetterCountry]; exists {
				tField.ID = element.ID
				tField.Name = element.Role
				tField.Name = utils.GetStringValue(record.NationalID, "")
			}
		}

		fields = append(fields, tField)
	}

	return fields
}
