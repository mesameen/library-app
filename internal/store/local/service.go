package local

import "github.com/test/library-app/internal/model"

type LocalStore struct{}

func (l *LocalStore) GetBook() (*model.Book, error) {
	return nil, nil
}
