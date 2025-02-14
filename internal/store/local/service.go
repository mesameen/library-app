package local

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

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
	titleLower := strings.ToLower(title)
	// retireving it from store
	book, ok := l.books[titleLower]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with tile '%s' isn't presents", title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return book, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	return book, nil
}

// AddLoan adds the loan details to store
func (l *LocalStore) AddLoan(ctx context.Context, det *model.LoanDetails) error {
	loanID := GetUniqueIncrementedID()
	l.loans[loanID] = det
	return nil
}

var uniqueID int32

func GetUniqueIncrementedID() int {
	// incrementing uniqueID
	atomic.AddInt32(&uniqueID, 1)
	return int(uniqueID)
}

// Close clears the memory
func (l *LocalStore) Close() error {
	logger.Infof("clearing up local store")
	// clearing it up local store
	l.books = nil
	l.loans = nil
	return nil
}
