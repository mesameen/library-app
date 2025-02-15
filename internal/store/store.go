package store

import (
	"context"
	"fmt"

	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/constants"
	"github.com/test/library-app/internal/model"
	"github.com/test/library-app/internal/store/local"
	"github.com/test/library-app/internal/store/postgres"
)

type Store interface {
	// GetBookDetails retreves book details from store
	GetBookDetails(ctx context.Context, title string) (*model.BookDetails, error)
	// GetAllBookDetails retreves book details from store
	GetAllBookDetails(ctx context.Context) ([]*model.BookDetails, error)
	// AddLoan adds the loan details to store
	AddLoan(ctx context.Context, det *model.LoanDetails) (int, error)
	// Extends the loan
	ExtendLoan(ctx context.Context, loanID int) (*model.LoanDetails, error)
	// Retunrs a book
	ReturnBook(ctx context.Context, loanID int) error
	Close() error
}

func NewStore() (Store, error) {
	switch config.CommonConfig.StoreType {
	case constants.LocalStore:
		return local.InitLocalStore()
	case constants.PostgresStore:
		return postgres.InitPostgresStore()
	default:
		return nil, fmt.Errorf("unknown Store configured: %v", config.CommonConfig.StoreType)
	}
}
