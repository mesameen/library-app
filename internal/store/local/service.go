package local

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/test/library-app/internal/constants"
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
	l.rmu.RLock()
	defer l.rmu.RUnlock()
	books := make([]*model.BookDetails, 0)
	for _, book := range l.books {
		books = append(books, book)
	}
	return books, nil
}

func (l *LocalStore) GetAllLoans(ctx context.Context) ([]*model.LoanDetails, error) {
	l.rmu.RLock()
	defer l.rmu.RUnlock()
	loans := make([]*model.LoanDetails, 0)
	for _, loan := range l.loans {
		loans = append(loans, loan)
	}
	return loans, nil
}

// GetBookDetails retreves book details from store
func (l *LocalStore) GetBookDetails(ctx context.Context, title string) (*model.BookDetails, error) {
	l.rmu.RLock()
	defer l.rmu.RUnlock()
	// retireving it from store
	book, ok := l.books[strings.ToLower(title)]
	if !ok {
		// If requested title isn't presents returning error with info,
		// err := fmt.Errorf("book with title '%s' isn't presents", title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return nil, fmt.Errorf("book with title '%s' isn't presents. %w", title, model.ErrNotFound)
	}
	return book, nil
}

// AddLoan adds the loan details to store
func (l *LocalStore) AddLoan(ctx context.Context, det *model.LoanDetails) (int, error) {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	book, ok := l.books[strings.ToLower(det.Title)]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with title '%s' isn't presents", det.Title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return 0, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	// if available copies are zero returning the error
	if book.AvailableCopies == 0 {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with title '%s' are out of stock", det.Title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return 0, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	// getting unique id
	id := GetUniqueIncrementedID()
	det.ID = id
	// setting in to detailsshort
	l.loans[id] = det

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
	return id, nil
}

var uniqueID int32

func GetUniqueIncrementedID() int {
	// incrementing uniqueID
	atomic.AddInt32(&uniqueID, 1)
	return int(uniqueID)
}

// ExtendLoan by given value
func (l *LocalStore) ExtendLoan(ctx context.Context, loanID int) (*model.LoanDetails, error) {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	loan, ok := l.loans[loanID]
	if !ok {
		// If requested loan isn't presents returning error with info
		err := fmt.Errorf("loan %d isn't presents", loanID)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return nil, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	if loan.Status == constants.Closed {
		logger.Errorf("requested loan: %d already closed", loanID)
		return nil, fmt.Errorf("requested loan: %d already closed", loanID)
	}
	returnTime := time.Unix(loan.ReturnDate, 0)
	// extending 3 weeks
	loan.ReturnDate = returnTime.Add(24 * 7 * 3 * time.Hour).Unix()
	logger.Infof("Loan extended for book title: %s", loan.Title)
	return loan, nil
}

// ExtendLoan by given value
func (l *LocalStore) ReturnBook(ctx context.Context, loanID int) (*model.LoanDetails, error) {
	l.rmu.Lock()
	defer l.rmu.Unlock()
	loan, ok := l.loans[loanID]
	if !ok {
		// If requested loan isn't presents returning error with info
		err := fmt.Errorf("loan %d isn't presents", loanID)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return nil, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	if loan.Status == constants.Closed {
		logger.Errorf("requested loan: %d already closed", loanID)
		return nil, fmt.Errorf("requested loan: %d already closed", loanID)
	}
	// reducing one from the avalilablecopies of the title
	bookDet, ok := l.books[strings.ToLower(loan.Title)]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with title '%s' isn't presents", loan.Title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return nil, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	// reducing one from available copies
	bookDet.AvailableCopies += 1
	logger.Infof("Title: %s is returned", loan.Title)

	// removing the loan from cache since book is returned
	loan.Status = constants.Closed
	logger.Infof("title: %s returned", loan.Title)
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
