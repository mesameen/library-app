package model

import (
	"errors"
)

// BookDetail represents book details
type BookDetails struct {
	Title           string `json:"title" example:"alchemist"`     // Unique Identifier for the book
	AvailableCopies int    `json:"available_copies" example:"10"` // No of available copies of the book that can be loaned
}

// LoanDetails represents loan of the book
type LoanDetails struct {
	ID             int    `json:"id"`               // auto generated at the backend
	NameOfBorrower string `json:"name_of_borrower"` // Name of borrower
	Title          string `json:"title"`            // title of the book
	LoanDate       int64  `json:"loan_date"`        // Date when the book was borrowed, unix epoch format. relavant for api calls
	ReturnDate     int64  `json:"return_date"`      // Date when the book should be returned, unix epoch format. relavant for api calls
	Status         string `json:"status"`           // active | closed
	// TODO: adding extra fields for additional functionality
	// RentPerDay     int    `json:"cost_per_day"`
}

// LoanDetails request
type LoanRequest struct {
	// binding: required
	NameOfBorrower string `json:"name_of_borrower" example:"john"` // Name of borrower
	Title          string `json:"title" example:"alchemist"`       // title of the book
}

// Custom Errors
var (
	ErrNotFound = errors.New("not found")
)

// CustomError
type CustomError struct {
	Error   string `json:"error"`
	Details string `json:"details"`
	Code    int    `json:"code"`
}
