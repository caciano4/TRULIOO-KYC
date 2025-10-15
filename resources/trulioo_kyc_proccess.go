package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"trullio-kyc/config"
	"trullio-kyc/middleware"
	"trullio-kyc/models"
	"trullio-kyc/utils"
)

// extractReportsData extracts and counts reports from Trulioo response
func extractReportsData(clientDetails models.ClientDetailsResponse) (int, int, int, string) {
	amCount, wlCount, pepCount := 0, 0, 0
	allReports := make(map[string][]interface{})

	// Parse the response as generic interface to handle dynamic structure
	responseBytes, err := json.Marshal(clientDetails)
	if err != nil {
		return 0, 0, 0, ""
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(responseBytes, &responseData); err != nil {
		return 0, 0, 0, ""
	}

	if flowData, ok := responseData["flowData"].(map[string]interface{}); ok {
		for _, flow := range flowData {
			if flowMap, ok := flow.(map[string]interface{}); ok {
				if serviceData, ok := flowMap["serviceData"].([]interface{}); ok {
					for _, service := range serviceData {
						if serviceMap, ok := service.(map[string]interface{}); ok {
							if fullDetails, ok := serviceMap["fullServiceDetails"].(map[string]interface{}); ok {
								if record, ok := fullDetails["Record"].(map[string]interface{}); ok {
									if datasourceResults, ok := record["DatasourceResults"].([]interface{}); ok {
										for _, datasource := range datasourceResults {
											if dsMap, ok := datasource.(map[string]interface{}); ok {
												if fields, ok := dsMap["DatasourceFields"].([]interface{}); ok {
													for _, field := range fields {
														if fieldMap, ok := field.(map[string]interface{}); ok {
															if fieldMap["FieldName"] == "WatchlistHitDetails" {
																if data, ok := fieldMap["Data"].(map[string]interface{}); ok {
																	if amResults, ok := data["AM_results"].([]interface{}); ok && len(amResults) > 0 {
																		amCount += len(amResults)
																		allReports["AM_results"] = append(allReports["AM_results"], amResults...)
																	}
																	if wlResults, ok := data["WL_results"].([]interface{}); ok && len(wlResults) > 0 {
																		wlCount += len(wlResults)
																		allReports["WL_results"] = append(allReports["WL_results"], wlResults...)
																	}
																	if pepResults, ok := data["PEP_results"].([]interface{}); ok && len(pepResults) > 0 {
																		pepCount += len(pepResults)
																		allReports["PEP_results"] = append(allReports["PEP_results"], pepResults...)
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	reportsJson := ""
	if len(allReports) > 0 {
		if jsonData, err := json.Marshal(allReports); err == nil {
			reportsJson = string(jsonData)
		}
	}

	return amCount, wlCount, pepCount, reportsJson
}

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

type TruliooSession struct {
	XHfSession  string
	BearerToken string
	JobID       string
}

type Fields []TField

// getTruliooBaseURL returns the base URL based on environment
func getTruliooBaseURL() string {
	env := config.GetEnv("TRULIOO_ENV", "test")
	config.AppLogger.Printf("Using Trulioo environment: %s", env)
	if env == "prod" {
		return "https://api.workflow.prod.trulioo.com/interpreter-v2"
	}
	return "https://api.workflow.prod.trulioo.com/interpreter-v2/test"
}

// getTruliooExportURL returns the export URL based on environment
func getTruliooExportURL() string {
	env := config.GetEnv("TRULIOO_ENV", "test")
	if env == "prod" {
		return "https://api.workflow.prod.trulioo.com/export/v2"
	}
	return "https://api.workflow.prod.trulioo.com/export/test/v2"
}

// getTruliooAuthURL returns the auth URL (same for both environments)
func getTruliooAuthURL() string {
	return "https://auth-api.trulioo.com/connect/token"
}

func HandleCatchKYCById(r *http.Request) ([]models.Record, error) {
	// Track Log
	config.AppLogger.Print("Get KYC RECORDS")

	// Start variables
	param := r.Context().Value(middleware.ParamsKey).(string)
	var records []models.Record

	// Connect and Close database
	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	// Preparing query to fetch KYC by package_file_id (only pending records)
	query := `
			SELECT
				id, package_file_id, package_name, upload_by_id, client_reference_id, transfer_agent_responsible,
				type_of_transfer, email, user_id, first_name, middle_name, last_name, date_of_birth_day,
				personal_phone_number, street_address, city, postal, letter_state, letter_country,
				national_id, request, response, notes, match, complete_kyc, created_at, updated_at, deleted_at
			FROM public.document_records
			WHERE package_file_id = $1 AND complete_kyc = false
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

func HandleProcessAllKyc(w http.ResponseWriter, r *http.Request, record models.Record) error {
	// Create isolated session for this job
	session := &TruliooSession{
		JobID: fmt.Sprintf("kyc-%d-%d", record.Id, time.Now().Unix()),
	}

	config.AppLogger.Printf("🔐 Starting isolated KYC session %s for record ID %d", session.JobID, record.Id)

	// Step 1: Init and catch field Ids
	fields, err := truliooInit(w, record)
	if err != nil {
		return fmt.Errorf("failed to initialize Trulioo fields for record Id %v: %w", record.Id, err)
	}

	// Step 2: Send body with Ids, and store the request sent
	err = truliooBodySubmit(fields, record, session)
	if err != nil {
		return fmt.Errorf("failed to submit Trulioo body for record Id %v: %w", record.Id, err)
	}
	config.AppLogger.Printf("🔑 Session %s acquired XHfSession: %s", session.JobID, session.XHfSession)

	// Step 3: Retrieve Bearer Token
	err = truliooGenerateBearerToken(record, session)
	if err != nil {
		return fmt.Errorf("failed to generate Bearer token for record Id %v: %w", record.Id, err)
	}
	config.AppLogger.Printf("🎫 Session %s acquired Bearer Token: %s...", session.JobID, session.BearerToken[:20])

	// Step 4: Match API Trulioo
	err = truliooDetailsFromClient(w, r, record, session)
	if err != nil {
		return fmt.Errorf("failed to retrieve details from Trulioo client for record Id %v: %w", record.Id, err)
	}

	// If all steps succeed, return nil
	config.AppLogger.Printf("✅ Session %s completed successfully for record ID %d", session.JobID, record.Id)
	return nil
}

// step 4
func truliooDetailsFromClient(w http.ResponseWriter, r *http.Request, record models.Record, session *TruliooSession) error {
	config.AppLogger.Print("TRULIOO DETAILS FROM CLIENT: STEP 4")
	var completed bool = false
	var clientDetails models.ClientDetailsResponse
	var request Req
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)
	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	request.URL = fmt.Sprintf("%s/query/client/%s?includeFullServiceDetails=true", getTruliooExportURL(), session.XHfSession)

	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", session.BearerToken))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	if err := json.Unmarshal(body, &clientDetails); err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	completed = true
	config.LogResponseTrulio(4, userName, clientDetails, "response")

	// Extract reports data
	amCount, wlCount, pepCount, reportsData := extractReportsData(clientDetails)

	// Update record with match data and mark as completed
	query := `UPDATE document_records
				SET match = $1, complete_kyc = $3, response = $4, 
					am_reports_count = $5, wl_reports_count = $6, pep_reports_count = $7, reports_data = $8,
					updated_at = NOW()
			 WHERE id = $2`

	flowDataJson, err := json.Marshal(clientDetails.FlowData)
	if err != nil {
		config.AppLogger.Print(fmt.Sprintf("Error serializing FlowData: %v", err))
		return err
	}

	_, err = db.Exec(query, "true", record.Id, completed, string(flowDataJson), amCount, wlCount, pepCount, reportsData)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	config.AppLogger.Printf("✅ Record %d updated with reports: AM=%d, WL=%d, PEP=%d", record.Id, amCount, wlCount, pepCount)
	return nil
}

// step 3
func truliooGenerateBearerToken(record models.Record, session *TruliooSession) error {
	config.AppLogger.Print("TRULIOO GENERATE BEARER TOKEN: STEP 3")
	var request Req
	var bearerTokenResponse models.BearerTokenReponse
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)

	//Using that type of body (x-www-form-urlencoded) because it's required as an oauth2 api
	request.URL = getTruliooAuthURL()
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
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	if err := json.Unmarshal(body, &bearerTokenResponse); err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	config.LogResponseTrulio(3, userName, bearerTokenResponse, "response")

	session.BearerToken = bearerTokenResponse.AccessToken

	return nil
}

// step 2
func truliooBodySubmit(fields Fields, record models.Record, session *TruliooSession) error {
	config.AppLogger.Print("TRULIOO SUBMIT: STEP 2")
	var request Req
	var truliooBodySubmitResponse models.DirectSubmitResponse
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)
	if request.Body == nil {
		request.Body = make(map[string]interface{})
	}

	request.FlowId = config.GetEnv("FLOW_ID", "")
	request.URL = fmt.Sprintf("%s/submit/%s", getTruliooBaseURL(), request.FlowId)

	for _, field := range fields {
		if field.ID != "" && field.TValue != "" {
			request.Body[field.ID] = field.TValue
		}
	}

	// Convert body map json
	bodyJson, err := json.Marshal(request.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	config.LogResponseTrulio(2, userName, request.Body, "request")

	req, err := http.NewRequest("POST", request.URL, bytes.NewBuffer(bodyJson))
	if err != nil {
		config.AppLogger.Printf(err.Error())
		return err
	}

	// Setting Header content type to json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hf-retry-on-pending", "true")

	// Submiting KYC request step 2
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}
	defer res.Body.Close()

	//GET RESPONSE FROM STEP 2o
	body, err := io.ReadAll(res.Body)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	// Decoding json HTTP
	if err := json.Unmarshal(body, &truliooBodySubmitResponse); err != nil {
		config.AppLogger.Print(err.Error())
		return err
	}

	config.LogResponseTrulio(2, userName, truliooBodySubmitResponse, "response")

	//! TODO Add resposne return from trulioo in DB
	//! TODO Update ROW with response text
	// truliooBodySubmitResponse.Text

	// GETTING XHFSESSION
	session.XHfSession = res.Header.Get("x-hf-session")
	return err
}

// ! Step 1
func truliooInit(w http.ResponseWriter, record models.Record) (Fields, error) {
	config.AppLogger.Print("INIT TRULIOO REQUEST: STEP 1")

	// Variables
	var request Req
	var initTrulioo models.InitTrulioo
	var fields Fields
	var tField TField
	userName := fmt.Sprintf("%s_%s", *record.FirstName, *record.LastName)

	// Preparing Req struct
	request.FlowId = config.GetEnv("FLOW_ID", "")
	request.URL = fmt.Sprintf("%s/flow/%s", getTruliooBaseURL(), request.FlowId)

	// Instance new Request
	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		config.AppLogger.Print(err.Error())
		return Fields{}, nil
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
		return Fields{}, nil
	}

	// Decoding Json response (se aplicável)
	if err := json.Unmarshal(body, &initTrulioo); err != nil {
		http.Error(w, "Error to decode json", http.StatusInternalServerError)
		config.AppLogger.Print(err.Error())
		return Fields{}, nil
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

	return fields, nil
}
