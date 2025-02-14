package local

import (
	"fmt"
	"strings"

	"github.com/test/library-app/internal/logger"
	"github.com/test/library-app/internal/model"
)

type LocalStore struct {
	store map[string]*model.BookDetail // stores the Books key as book tiltle
}

// GetBookDetails retreves book details from store
func (l *LocalStore) GetBookDetails(title string) (*model.BookDetail, error) {
	titleLower := strings.ToLower(title)
	// retireving it from store
	book, ok := l.store[titleLower]
	if !ok {
		// If requested title isn't presents returning error with info,
		err := fmt.Errorf("book with tile '%s' isn't presents", title)
		// wrapping with NotFound error to identify the error type by caller or middleware
		return book, fmt.Errorf("%v %w", err, model.ErrNotFound)
	}
	return book, nil
}

func (l *LocalStore) Close() error {
	logger.Infof("clearing up local store")
	// clearing it up local store
	l.store = nil
	return nil
}
