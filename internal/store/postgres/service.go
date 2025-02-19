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
	query := fmt.Sprintf(`SELECT 
		title, 
		available_copies 
		FROM %s
	`, config.PostgresConfig.BooksTableName)
	rows, err := p.DB.Query(ctx, query)
	if err != nil {
		logger.Errorf("Failed to fetch books. Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to fetch books. %w", model.ErrNotFound)
		}
		return nil, err
	}
	defer rows.Close()
	books := make([]*model.BookDetails, 0)
	for rows.Next() {
		var book model.BookDetails
		if err := rows.Scan(&book.Title, &book.AvailableCopies); err != nil {
			logger.Errorf("Failed to scan bookdetails fetched from DB. Error: %v", err)
			continue
		}
		books = append(books, &book)
	}

	return books, nil
}

// GetAllLoans retreves all loan details from store
func (p *PostgresDB) GetAllLoans(ctx context.Context) ([]*model.LoanDetails, error) {
	query := fmt.Sprintf(`SELECT 
		id,
		title, 
		name_of_borrower,
		loan_date,
		return_date 
		FROM %s
	`, config.PostgresConfig.LoansTableName)
	rows, err := p.DB.Query(ctx, query)
	if err != nil {
		logger.Errorf("Failed to fetch loans. Error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to fetch books. %w", model.ErrNotFound)
		}
		return nil, err
	}
	defer rows.Close()
	loans := make([]*model.LoanDetails, 0)
	for rows.Next() {
		var loan model.LoanDetails
		var loanDate time.Time
		var returnDate time.Time
		if err := rows.Scan(&loan.ID, &loan.Title, &loan.NameOfBorrower, &loanDate, &returnDate); err != nil {
			logger.Errorf("Failed to scan bookdetails fetched from DB. Error: %v", err)
			continue
		}
		loan.LoanDate = loanDate.Unix()
		loan.ReturnDate = returnDate.Unix()
		loans = append(loans, &loan)
	}

	return loans, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("failed to find the title: %s. %w", det.Title, model.ErrNotFound)
		}
		return 0, err
	}
	// if available copies are zero returning the error
	if avalilableCopies == 0 {
		logger.Errorf("not enough copies of requested title %v", det.Title)
		return 0, fmt.Errorf("not enough copies of requested title %v. %w", det.Title, model.ErrNotFound)
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
	// updating the return date
	query := fmt.Sprintf(`UPDATE
	%s SET return_date=return_date + interval '3 weeks'
	WHERE id=$1
	`, config.PostgresConfig.LoansTableName)
	_, err = tx.Exec(ctx, query, loanID)
	if err != nil {
		logger.Errorf("Failed to execute update query for extending loan. Error: %v", err)
		// if update failed with doesn't exists returning loan not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find loan: %d. %w", loanID, model.ErrNotFound)
		}
		return nil, err
	}
	// fetching the updated return date
	query = fmt.Sprintf(`SELECT
		return_date 
	FROM %s 
		WHERE id=$1 
	`, config.PostgresConfig.LoansTableName)
	var returnDate time.Time
	err = tx.QueryRow(ctx, query, loanID).Scan(&returnDate)
	if err != nil {
		logger.Errorf("failed to find a requested loan: %d to extend", loanID)
		return nil, fmt.Errorf("failed to find a requested loan: %d to extend. %w", loanID, model.ErrNotFound)
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
	// fetching title from loan
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

	// deleting the loan
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
