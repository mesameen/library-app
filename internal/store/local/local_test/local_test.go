package localtest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test/library-app/internal/store/local"
)

func TestInitLocalStore(t *testing.T) {
	ctx := context.Background()
	store, err := local.InitLocalStore()
	assert.Nil(t, err)
	assert.NotNil(t, store)

	books, err := store.GetAllBookDetails(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, store)
	assert.GreaterOrEqual(t, len(books), 0)

	// success case
	book, err := store.GetBookDetails(ctx, "book_1")
	assert.Nil(t, err)
	assert.NotNil(t, book)

	// failure case
	book, err = store.GetBookDetails(ctx, "book_xyz")
	assert.NotNil(t, err)
	assert.Nil(t, book)
}
