package local

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/test/library-app/internal/logger"
	"github.com/test/library-app/internal/model"
)

// making the members of store as private to avoid updating from elsewhere other than the allowed functions
type LocalStore struct {
	rmu   sync.RWMutex
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
	l.rmu.RLock()
	defer l.rmu.RUnlock()
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
	l.rmu.Lock()
	defer l.rmu.Unlock()
	l.loans[det.ID] = det

	// reducing one from the avalilablecopies of the title
	bookDet, ok := l.books[strings.ToLower(det.Title)]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with title '%s' isn't presents", det.Title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return 0, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	// reducing one from available copies
	bookDet.AvailableCopies -= 1

	logger.Infof("Loan entry added for book title: %s", det.Title)
	return det.ID, nil
}

// ExtendLoan by given value
func (l *LocalStore) ExtendLoan(ctx context.Context, loanID int) (*model.LoanDetails, error) {
	l.rmu.Lock()
	defer l.rmu.Unlock()
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

// ExtendLoan by given value
func (l *LocalStore) ReturnBook(ctx context.Context, loanID int) error {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	loan, ok := l.loans[loanID]
	if !ok {
		// If requested loan isn't presents returning error with info
		err := fmt.Errorf("loan %d isn't presents", loanID)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	// reducing one from the avalilablecopies of the title
	bookDet, ok := l.books[strings.ToLower(loan.Title)]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with title '%s' isn't presents", loan.Title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	// reducing one from available copies
	bookDet.AvailableCopies += 1
	logger.Infof("Title: %s is returned", loan.Title)

	// removing the loan from cache since book is returned
	delete(l.loans, loanID)
	logger.Infof("title: %s returned", loan.Title)
	return nil
}

// Close clears the memory
func (l *LocalStore) Close() error {
	logger.Infof("clearing up local store")
	// clearing it up local store
	l.books = nil
	l.loans = nil
	return nil
}
