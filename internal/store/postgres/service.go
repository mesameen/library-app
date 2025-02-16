package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
	tx, err := p.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Errorf("failed to begin transaction. Error: %v", err)
		return 0, err
	}
	defer tx.Rollback(ctx)
	query := fmt.Sprintf(`SELECT available_copies FROM %s WHERE LOWER(title)=LOWER($1)`, config.PostgresConfig.BooksTableName)
	var avalilableCopies int
	err = tx.QueryRow(ctx, query, det.Title).Scan(&avalilableCopies)
	if err != nil {
		logger.Errorf("failed to fetch requested title from books table. Error: %v", err)
		return 0, err
	}
	// if available copies are zero returning the error
	if avalilableCopies == 0 {
		logger.Errorf("not enough copies of requested title %v", det.Title)
		return 0, fmt.Errorf("not enough copies of requested title %v", det.Title)
	}
	// id := GetUniqueIncrementedID()
	lastInsertId := 0
	// inserting in to loans table
	query = fmt.Sprintf(`INSERT
		INTO %s
		(title, name_of_borrower, return_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, config.PostgresConfig.LoansTableName)
	err = tx.QueryRow(ctx, query, det.Title, det.NameOfBorrower, time.Unix(det.ReturnDate, 0)).Scan(&lastInsertId)
	if err != nil {
		logger.Errorf("failed to insert into loan. Error: %v", err)
		return 0, err
	}
	det.ID = lastInsertId

	// updating the available copies
	query = fmt.Sprintf(`UPDATE
		%s SET available_copies=$1 WHERE LOWER(title)=LOWER($2)
	`, config.PostgresConfig.BooksTableName)
	_, err = tx.Exec(ctx, query, avalilableCopies-1, det.Title)
	if err != nil {
		logger.Errorf("failed to update avaialble_copies count in to books. Error: %v", err)
		return 0, err
	}
	// committing the transaction after all db actions completed successfully
	if err = tx.Commit(ctx); err != nil {
		logger.Errorf("failed to commit transaction. Error: %v", err)
		return 0, err
	}
	return det.ID, nil
}

// Extends the loan
func (p *PostgresDB) ExtendLoan(ctx context.Context, loanID int) (*model.LoanDetails, error) {
	tx, err := p.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Errorf("Failed to begin transaction. Error: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	query := fmt.Sprintf(`UPDATE
	%s SET return_date=return_date + interval '3 weeks'
	WHERE id=$1
	`, config.PostgresConfig.LoansTableName)
	_, err = tx.Exec(ctx, query, loanID)
	if err != nil {
		logger.Errorf("Failed to execute update query for extending loan. Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find loan: %d. %w", loanID, model.ErrNotFound)
		}
		return nil, err
	}
	query = fmt.Sprintf(`SELECT
		return_date 
	FROM %s 
		WHERE id=$1 
	`, config.PostgresConfig.LoansTableName)
	var returnDate time.Time
	err = tx.QueryRow(ctx, query, loanID).Scan(&returnDate)
	if err != nil {
		logger.Errorf("failed to find updated loan: %d", loanID)
		return nil, fmt.Errorf("failed to find updated loan: %d. %w", loanID, model.ErrNotFound)
	}
	if err = tx.Commit(ctx); err != nil {
		logger.Errorf("Failed to commit transaction of extending loan. Error: %v", err)
		return nil, err
	}
	return &model.LoanDetails{
		ReturnDate: returnDate.Unix(),
	}, nil
}

// Retunrs a book
func (p *PostgresDB) ReturnBook(ctx context.Context, loanID int) error {
	tx, err := p.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Errorf("Failed to begin transaction. Error: %v", err)
		return err
	}
	defer tx.Rollback(ctx)
	var title string
	query := fmt.Sprintf(`SELECT
		title
	FROM
		%s
		WHERE id=$1
	`, config.PostgresConfig.LoansTableName)
	err = tx.QueryRow(ctx, query, loanID).Scan(&title)
	if err != nil {
		logger.Errorf("Failed to execute get loan. Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to find loan: %d. %w", loanID, model.ErrNotFound)
		}
		return err
	}

	query = fmt.Sprintf(`DELETE
		FROM
		%s
		WHERE id=$1
	`, config.PostgresConfig.LoansTableName)
	_, err = tx.Exec(ctx, query, loanID)
	if err != nil {
		logger.Errorf("Failed to execute update query for extending loan. Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to find loan: %d. %w", loanID, model.ErrNotFound)
		}
		return err
	}
	query = fmt.Sprintf(`UPDATE
		%s SET available_copies=available_copies+1
	WHERE 
		LOWER(title)=LOWER($1)
	`,
		config.PostgresConfig.BooksTableName)
	_, err = tx.Exec(ctx, query, title)
	if err != nil {
		logger.Errorf("Failed to update  query for extending loan. Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to find book: %d. %w", loanID, model.ErrNotFound)
		}
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		logger.Errorf("Failed to commit transaction of returning a book. Error: %v", err)
		return err
	}
	return nil
}

func (p *PostgresDB) Close() error {
	logger.Infof("Closing the postgress connection pool")
	p.DB.Close()
	return nil
}
