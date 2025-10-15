package models

import "time"

type KYCReportData struct {
	// Header Info
	ClientName    string `json:"client_name"`
	GeneratedDate string `json:"generated_date"`

	// Application Details
	ApplicationID     string `json:"application_id"`
	ClientReferenceID string `json:"client_reference_id"`
	FullName          string `json:"full_name"`
	FullAddress       string `json:"full_address"`
	DateOfBirth       string `json:"date_of_birth"`
	NationalID        string `json:"national_id"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`

	// Match Status
	NameMatch        string `json:"name_match"`
	NameMatchClass   string `json:"name_match_class"`
	AddressMatch     string `json:"address_match"`
	AddressMatchClass string `json:"address_match_class"`
	DOBMatch         string `json:"dob_match"`
	DOBMatchClass    string `json:"dob_match_class"`
	IDMatch          string `json:"id_match"`
	IDMatchClass     string `json:"id_match_class"`

	// Screening Results
	WatchlistHits      int `json:"watchlist_hits"`
	AdverseMediaHits   int `json:"adverse_media_hits"`
	PEPHits           int `json:"pep_hits"`

	// Processing Steps
	ProcessingSteps []ProcessingStep `json:"processing_steps"`

	// Overall Status
	Status string `json:"status"`

	// Additional Info
	TransactionID string `json:"transaction_id"`
	ProcessedDate string `json:"processed_date"`
}

type ProcessingStep struct {
	Step    string `json:"step"`
	Service string `json:"service"`
	Status  string `json:"status"`
	Match   string `json:"match"`
}

// Helper function to convert database record to PDF report data
func ConvertToPDFReport(record Record, responseData ClientDetailsResponse) KYCReportData {
	report := KYCReportData{
		GeneratedDate: time.Now().Format("2006-01-02"),
		DateOfBirth:   record.DateOfBirthDay.Format("2006-01-02"),
	}

	// Set ApplicationID
	if record.PackageID != nil && *record.PackageID != "" {
		report.ApplicationID = *record.PackageID
	} else {
		report.ApplicationID = "N/A"
	}

	// Set ClientReferenceID
	if record.ClientReferenceID != nil && *record.ClientReferenceID != "" {
		report.ClientReferenceID = *record.ClientReferenceID
	} else {
		report.ClientReferenceID = "N/A"
	}

	// Set Phone
	if record.PersonalPhoneNumber != nil && *record.PersonalPhoneNumber != "" {
		report.Phone = *record.PersonalPhoneNumber
	} else {
		report.Phone = "N/A"
	}

	// Set Email
	if record.Email != nil && *record.Email != "" {
		report.Email = *record.Email
	} else {
		report.Email = "N/A"
	}

	// Set NationalID
	if record.NationalID != nil && *record.NationalID != "" {
		report.NationalID = *record.NationalID
	} else {
		report.NationalID = "N/A"
	}

	// Build full name
	firstName := ""
	lastName := ""
	middleName := ""

	if record.FirstName != nil {
		firstName = *record.FirstName
	}
	if record.LastName != nil {
		lastName = *record.LastName
	}
	if record.MiddleName != nil {
		middleName = *record.MiddleName
	}

	if middleName != "" {
		report.FullName = firstName + " " + middleName + " " + lastName
	} else {
		report.FullName = firstName + " " + lastName
	}
	report.ClientName = report.FullName

	// Build full address
	address := ""
	city := ""
	state := ""
	postal := ""
	country := ""

	if record.StreetAddress != nil {
		address = *record.StreetAddress
	}
	if record.City != nil {
		city = *record.City
	}
	if record.LetterState != nil {
		state = *record.LetterState
	}
	if record.Postal != nil {
		postal = *record.Postal
	}
	if record.LetterCountry != nil {
		country = *record.LetterCountry
	}

	report.FullAddress = address + ", " + city + ", " + state + ", " + postal + ", " + country

	// Extract match data from response
	report.extractMatchData(responseData)

	// Extract screening results
	report.extractScreeningResults(responseData)

	// Extract processing steps
	report.extractProcessingSteps(responseData)

	// Set overall status
	report.Status = responseData.Status

	return report
}

func (r *KYCReportData) extractMatchData(responseData ClientDetailsResponse) {
	// Initialize with default values
	r.NameMatch = "No"
	r.NameMatchClass = "match-no"
	r.AddressMatch = "No"
	r.AddressMatchClass = "match-no"
	r.DOBMatch = "No"
	r.DOBMatchClass = "match-no"
	r.IDMatch = "No"
	r.IDMatchClass = "match-no"

	// For now, use simplified extraction
	// TODO: Implement proper FlowData parsing when structure is stable
	if responseData.Status == "ACCEPTED" {
		r.NameMatch = "Yes"
		r.NameMatchClass = "match-yes"
		r.AddressMatch = "Yes"
		r.AddressMatchClass = "match-yes"
		r.DOBMatch = "Yes"
		r.DOBMatchClass = "match-yes"
		r.IDMatch = "Yes"
		r.IDMatchClass = "match-yes"
	}
}

func (r *KYCReportData) extractScreeningResults(responseData ClientDetailsResponse) {
	// Initialize defaults
	r.WatchlistHits = 0
	r.AdverseMediaHits = 0
	r.PEPHits = 0

	// TODO: Implement proper screening results extraction
	// For now, use placeholder values
}

func (r *KYCReportData) extractProcessingSteps(responseData ClientDetailsResponse) {
	r.ProcessingSteps = []ProcessingStep{}

	// TODO: Implement proper processing steps extraction
	// For now, add placeholder steps
	r.ProcessingSteps = append(r.ProcessingSteps, ProcessingStep{
		Step:    "Step 1",
		Service: "Person Match",
		Status:  "COMPLETED",
		Match:   "Yes",
	})

	r.ProcessingSteps = append(r.ProcessingSteps, ProcessingStep{
		Step:    "Step 2",
		Service: "Watchlist Screening",
		Status:  "COMPLETED",
		Match:   "No",
	})

	// Set placeholder values
	r.TransactionID = "sample-transaction-id"
	r.ProcessedDate = time.Now().Format("2006-01-02 15:04:05")
}