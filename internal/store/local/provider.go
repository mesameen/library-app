package local

import "github.com/test/library-app/internal/model"

var localStore = make(map[string]*model.Book)

func InitLocalStore() (*LocalStore, error) {
	return nil, nil
}
