package local

import (
	"strings"

	"github.com/test/library-app/internal/model"
)

// InitLocalStore initializes with some book details
func InitLocalStore() (*LocalStore, error) {
	books := []*model.BookDetails{
		{
			Title:           "Alchemist",
			AvailableCopies: 3,
		},
		{
			Title:           "Atomic Habbits",
			AvailableCopies: 10,
		},
		{
			Title:           "Sapiens",
			AvailableCopies: 10,
		},
		{
			Title:           "Mocking Bird",
			AvailableCopies: 10,
		},
		{
			Title:           "Animal Farm",
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
