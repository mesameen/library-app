package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/test/library-app/internal/config"
	"github.com/test/library-app/internal/logger"
	"github.com/test/library-app/internal/model"
)

type PostgresDB struct {
	DB *pgxpool.Pool
}

// GetBookDetails retreves book details from store
func (p *PostgresDB) GetBookDetails(ctx context.Context, title string) (*model.BookDetails, error) {
	query := fmt.Sprintf(`SELECT 
		title, 
		available_copies 
		FROM %s
		WHERE LOWER(title)=LOWER($1)
	`, config.PostgresConfig.BooksTableName)
	row := p.DB.QueryRow(ctx, query, title)
	var bookTitle string
	var avalilableCopies int
	err := row.Scan(&bookTitle, &avalilableCopies)
	if err != nil {
		logger.Errorf("Failed to scan the requested title: %s. Error: %v", title, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find the title: %s. %w", title, model.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find the title: %s", title)
	}
	return &model.BookDetails{
		Title:           bookTitle,
		AvailableCopies: avalilableCopies,
	}, nil
}

// GetAllBookDetails retreves book details from store
func (p *PostgresDB) GetAllBookDetails(ctx context.Context) ([]*model.BookDetails, error) {
	return nil, nil
}

// AddLoan adds the loan details to store
func (p *PostgresDB) AddLoan(ctx context.Context, det *model.LoanDetails) (int, error) {
	return 0, nil
}

// Extends the loan
func (p *PostgresDB) ExtendLoan(ctx context.Context, loanID int) (*model.LoanDetails, error) {
	return nil, nil
}

// Retunrs a book
func (p *PostgresDB) ReturnBook(ctx context.Context, loanID int) error {
	return nil
}

func (p *PostgresDB) Close() error {
	logger.Infof("Closing the postgress connection pool")
	p.DB.Close()
	return nil
}
