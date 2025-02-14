package model

import (
	"errors"
	"time"
)

// BookDetail represents book details
type BookDetail struct {
	Title           string // Unique Identifier for the book
	AvailableCopies int    // No of available copies of the book that can be loaned
}

// LoanDetails represents loan of the book
type LoanDetails struct {
	NameOfBorrower string    // Name of borrower
	LoanDate       time.Time // Date when the book was borrowed
	ReturnDate     time.Time // Date when the book should be returned
}

// Custom Errors
var (
	ErrNotFound = errors.New("not found")
)
