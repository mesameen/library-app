package local

import "github.com/test/library-app/internal/model"

func InitLocalStore() (*LocalStore, error) {
	var localStore = make(map[int]*model.BookDetail)
	localStore[1] = &model.BookDetail{
		ID:              1,
		Title:           "Book 1",
		AvailableCopies: 10,
	}
	localStore[2] = &model.BookDetail{
		ID:              2,
		Title:           "Book 2",
		AvailableCopies: 10,
	}
	localStore[3] = &model.BookDetail{
		ID:              3,
		Title:           "Book 3",
		AvailableCopies: 10,
	}
	localStore[4] = &model.BookDetail{
		ID:              4,
		Title:           "Book 4",
		AvailableCopies: 10,
	}
	localStore[5] = &model.BookDetail{
		ID:              5,
		Title:           "Book 5",
		AvailableCopies: 10,
	}
	return &LocalStore{
		store: localStore,
	}, nil
}
