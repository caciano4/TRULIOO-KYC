package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"trullio-kyc/config"
	"trullio-kyc/middleware"
	"trullio-kyc/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateKYCPDF(w http.ResponseWriter, r *http.Request) {
	// Extract record ID from URL parameter
	recordIDStr := r.Context().Value(middleware.ParamsKey).(string)
	recordID, err := strconv.Atoi(recordIDStr)
	if err != nil {
		config.AppLogger.Printf("Invalid record ID: %s", recordIDStr)
		http.Error(w, "Invalid record ID", http.StatusBadRequest)
		return
	}

	config.AppLogger.Printf("Generating PDF report for record ID: %d", recordID)

	// Get record data from database
	record, err := getRecordByID(recordID)
	if err != nil {
		config.AppLogger.Printf("Error fetching record %d: %v", recordID, err)
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	// Parse response data if available
	var responseData models.ClientDetailsResponse
	if record.Response != nil && *record.Response != "" {
		err := json.Unmarshal([]byte(*record.Response), &responseData)
		if err != nil {
			config.AppLogger.Printf("Error parsing response data for record %d: %v", recordID, err)
			// Continue with empty response data
		}
	}

	// Convert to PDF report data
	reportData := models.ConvertToPDFReport(record, responseData)

	// Generate PDF directly
	pdfBytes, err := generatePDFReport(reportData)
	if err != nil {
		config.AppLogger.Printf("Error generating PDF for record %d: %v", recordID, err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Set response headers
	filename := fmt.Sprintf("kyc_report_%s_%d.pdf", reportData.ClientReferenceID, recordID)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	// Write PDF to response
	_, err = w.Write(pdfBytes)
	if err != nil {
		config.AppLogger.Printf("Error writing PDF response for record %d: %v", recordID, err)
	} else {
		config.AppLogger.Printf("Successfully generated PDF report for record %d", recordID)
	}
}

func getRecordByID(recordID int) (models.Record, error) {
	var record models.Record

	db := config.ConnectDB()
	defer config.CloseConnectionDB(db)

	query := `
		SELECT
			id, package_file_id, package_name, upload_by_id, client_reference_id, transfer_agent_responsible,
			type_of_transfer, email, user_id, first_name, middle_name, last_name, date_of_birth_day,
			personal_phone_number, street_address, city, postal, letter_state, letter_country,
			national_id, request, response, notes, match, complete_kyc, created_at, updated_at, deleted_at
		FROM public.document_records
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := db.QueryRow(query, recordID).Scan(
		&record.Id,
		&record.PackageID,
		&record.PackageName,
		&record.UploadById,
		&record.ClientReferenceID,
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
		&record.CompleteKYC,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.DeletedAt,
	)

	return record, err
}

func generatePDFReport(reportData models.KYCReportData) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add title
	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(0, 15, "KYC REPORT")
	pdf.Ln(20)

	// Add date
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Generated: "+reportData.GeneratedDate)
	pdf.Ln(15)

	// Application Details Section
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Application Details")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)

	// Create table-like structure
	addTableRow(pdf, "Application ID:", reportData.ApplicationID)
	addTableRow(pdf, "Client Reference:", reportData.ClientReferenceID)
	addTableRow(pdf, "Full Name:", reportData.FullName)
	addTableRow(pdf, "Address:", reportData.FullAddress)
	addTableRow(pdf, "Date of Birth:", reportData.DateOfBirth)
	addTableRow(pdf, "National ID/SSN:", reportData.NationalID)
	addTableRow(pdf, "Phone:", reportData.Phone)
	addTableRow(pdf, "Email:", reportData.Email)

	pdf.Ln(10)

	// Screening Details Section
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Screening Details")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	addTableRow(pdf, "Watch List Hits:", fmt.Sprintf("%d", reportData.WatchlistHits))
	addTableRow(pdf, "Adverse Media Hits:", fmt.Sprintf("%d", reportData.AdverseMediaHits))
	addTableRow(pdf, "PEP Hits:", fmt.Sprintf("%d", reportData.PEPHits))

	pdf.Ln(10)

	// Status Section
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Overall Status")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 14)
	if reportData.Status == "ACCEPTED" {
		pdf.SetTextColor(0, 128, 0) // Green
	} else {
		pdf.SetTextColor(255, 0, 0) // Red
	}
	pdf.Cell(0, 10, reportData.Status)
	pdf.SetTextColor(0, 0, 0) // Reset to black

	pdf.Ln(15)

	// Processing Steps
	if len(reportData.ProcessingSteps) > 0 {
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(0, 10, "Processing Steps")
		pdf.Ln(12)

		pdf.SetFont("Arial", "", 11)
		for _, step := range reportData.ProcessingSteps {
			addTableRow(pdf, step.Step+":", step.Service+" - "+step.Status+" (Match: "+step.Match+")")
		}
	}

	pdf.Ln(10)

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 8, "Generated by KORE KYC System - "+reportData.GeneratedDate)

	// Return PDF as bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error outputting PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func addTableRow(pdf *gofpdf.Fpdf, label, value string) {
	// Label column (30% width)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(60, 8, label, "1", 0, "L", false, 0, "")

	// Value column (70% width)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 8, value, "1", 1, "L", false, 0, "")
}