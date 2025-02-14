package local

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/test/library-app/internal/logger"
	"github.com/test/library-app/internal/model"
)

// making the members of store as private to avoid updating from elsewhere other than the allowed functions
type LocalStore struct {
	books map[string]*model.BookDetails // stores the Books key as book tiltle
	loans map[int]*model.LoanDetails    // stores the loans key as loan ID
}

func (l *LocalStore) GetAllBookDetails(ctx context.Context) ([]*model.BookDetails, error) {
	books := make([]*model.BookDetails, len(l.books))
	for _, book := range l.books {
		books = append(books, book)
	}
	return books, nil
}

// GetBookDetails retreves book details from store
func (l *LocalStore) GetBookDetails(ctx context.Context, title string) (*model.BookDetails, error) {
	// retireving it from store
	book, ok := l.books[strings.ToLower(title)]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with title '%s' isn't presents", title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return nil, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	return book, nil
}

// AddLoan adds the loan details to store
func (l *LocalStore) AddLoan(ctx context.Context, det *model.LoanDetails) (int, error) {
	l.loans[det.ID] = det

	// reducing one from the avalilablecopies of the title
	bookDet := l.books[strings.ToLower(det.Title)]
	// reducing one from available copies
	bookDet.AvailableCopies -= 1

	logger.Infof("Loan entry added for book title: %s", det.Title)
	return det.ID, nil
}

// ExtendLoan by given value
func (l *LocalStore) ExtendLoan(ctx context.Context, loanID int) (*model.LoanDetails, error) {
	loan, ok := l.loans[loanID]
	if !ok {
		// If requested loan isn't presents returning error with info
		err := fmt.Errorf("loan isn't %d isn't presents", loanID)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return nil, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	returnTime := time.Unix(loan.ReturnDate, 0)
	// extending 3 weeks
	loan.ReturnDate = returnTime.Add(24 * 7 * 3 * time.Hour).Unix()
	logger.Infof("Loan extended for book title: %s", loan.Title)
	return loan, nil
}

// Close clears the memory
func (l *LocalStore) Close() error {
	logger.Infof("clearing up local store")
	// clearing it up local store
	l.books = nil
	l.loans = nil
	return nil
}
