package localtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/test/library-app/internal/model"
	"github.com/test/library-app/internal/store/local"
)

var localStore *local.LocalStore
var err error
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	localStore, err = local.InitLocalStore()
	m.Run()
}

func TestInitLocalStore(t *testing.T) {
	assert.Nil(t, err)
	assert.NotNil(t, localStore)
}

func TestGetAllBookDetails(t *testing.T) {
	store, err := local.InitLocalStore()
	assert.Nil(t, err)
	assert.NotNil(t, store)

	books, err := store.GetAllBookDetails(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, store)
	assert.GreaterOrEqual(t, len(books), 0)
}

func TestGetBookDetails(t *testing.T) {
	// success case
	book, err := localStore.GetBookDetails(ctx, "book_1")
	assert.Nil(t, err)
	assert.NotNil(t, book)

	// failure case
	book, err = localStore.GetBookDetails(ctx, "book_xyz")
	assert.NotNil(t, err)
	assert.Nil(t, book)
}

func TestAddLoan(t *testing.T) {
	// success case
	loanID, err := localStore.AddLoan(ctx, &model.LoanDetails{
		ID:             1,
		NameOfBorrower: "test_user",
		Title:          "Book_1",
		LoanDate:       time.Now().Unix(),
		ReturnDate:     time.Now().Add(24 * time.Hour).Unix(),
	})
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, loanID, 0)

	// failure case
	loanID, err = localStore.AddLoan(ctx, &model.LoanDetails{
		ID:             2,
		NameOfBorrower: "test_user",
		Title:          "Book_100",
		LoanDate:       time.Now().Unix(),
		ReturnDate:     time.Now().Add(24 * time.Hour).Unix(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, loanID, 0)
}

func TestExtendLoan(t *testing.T) {
	// success case
	det, err := localStore.ExtendLoan(ctx, 1)
	assert.Nil(t, err)
	assert.NotNil(t, det)

	// failure case
	det, err = localStore.ExtendLoan(ctx, 10)
	assert.NotNil(t, err)
	assert.Nil(t, det)
}

func TestReturnBook(t *testing.T) {
	// success case
	err := localStore.ReturnBook(ctx, 1)
	assert.Nil(t, err)

	// failure case
	err = localStore.ReturnBook(ctx, 10)
	assert.NotNil(t, err)
}

func TestClose(t *testing.T) {
	err := localStore.Close()
	assert.Nil(t, err)
}
