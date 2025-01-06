package models

import "time"

type Record struct {
	Id                       int        `json:"id"` // New field
	PackageID                *string    `json:"package_id"`
	PackageName              *string    `json:"package_name"`
	UploadById               *string    `json:"upload_by_id"` // New field
	ClientReferenceID        *string    `json:"client_reference_id"`
	TransferAgentResponsible *string    `json:"transfer_agent_responsible"`
	TypeOfTransfer           *string    `json:"type_of_transfer"`
	Email                    *string    `json:"email"`
	UserID                   *string    `json:"user_id"`
	FirstName                *string    `json:"first_name"`
	MiddleName               *string    `json:"middle_name"`
	LastName                 *string    `json:"last_name"`
	DateOfBirthDay           time.Time  `json:"date_of_birth_day"`
	PersonalPhoneNumber      *string    `json:"personal_phone_number"`
	StreetAddress            *string    `json:"street_address"`
	City                     *string    `json:"city"`
	Postal                   *string    `json:"postal"`
	LetterState              *string    `json:"letter_state"`
	LetterCountry            *string    `json:"letter_country"`
	NationalID               *string    `json:"national_id"`
	Request                  *string    `json:"request"`
	Response                 *string    `json:"response"`
	Notes                    *string    `json:"notes"`
	Match                    *string    `json:"match"`
	CompleteKYC              bool       `json:"complete_kyc"` // New field
	VoterID                  *string    `json:"voter_id"`
	Passport                 *string    `json:"passport"`
	DriverLicence            *string    `json:"driver_license"`
	DriverLicenceVersion     *string    `json:"driver_license_version"`
	Suburb                   *string    `json:"suburb"`
	CreatedAt                time.Time  `json:"created_at"`           // New field
	UpdatedAt                time.Time  `json:"updated_at"`           // New field
	DeletedAt                *time.Time `json:"deleted_at,omitempty"` // New field
}
