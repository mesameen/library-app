package local

import (
	"github.com/test/library-app/internal/logger"
	"github.com/test/library-app/internal/model"
)

type LocalStore struct {
	store map[int]*model.BookDetail // stores the Books and info
}

func (l *LocalStore) GetBook() (*model.BookDetail, error) {
	return nil, nil
}

func (l *LocalStore) Close() error {
	logger.Infof("clearing up local store")
	// clearing it up local store
	l.store = nil
	return nil
}
