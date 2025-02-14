package local

import (
	"strings"

	"github.com/test/library-app/internal/model"
)

// InitLocalStore initializes with some book details
func InitLocalStore() (*LocalStore, error) {
	books := []*model.BookDetails{
		{
			Title:           "Book_1",
			AvailableCopies: 10,
		},
		{
			Title:           "Book_2",
			AvailableCopies: 10,
		},
		{
			Title:           "Book_3",
			AvailableCopies: 10,
		},
		{
			Title:           "Book_4",
			AvailableCopies: 10,
		},
		{
			Title:           "Book_5",
			AvailableCopies: 10,
		},
	}
	var localStore = make(map[string]*model.BookDetails)
	for _, book := range books {
		// lowering the title to keep it as key
		title := strings.ToLower(book.Title)
		localStore[title] = book
	}
	return &LocalStore{
		books: localStore,
		loans: make(map[int]*model.LoanDetails), // initializing the map
	}, nil
}
