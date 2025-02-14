package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/test/library-app/internal/model"
	"github.com/test/library-app/internal/store"
)

type Handler struct {
	repo store.Store
}

// Initializes requests handler
func NewHandler(s store.Store) *Handler {
	return &Handler{
		repo: s,
	}
}

// Test api
func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, John")
}

// GetBook retrieves the detail and available copies of a book title
func (h *Handler) GetBook(c *gin.Context) {
	title := c.Param("title")
	if len(title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is mandatory"})
	}
	det, err := h.repo.GetBookDetails(c, title)
	if err != nil {
		fmt.Println(errors.Is(err, model.ErrNotFound))
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			// unwrapping to send actual error
			err = errors.Unwrap(err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// rest of all errors falls under this category
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, det)
}

// BorrowBook borrows a book from store and returns the details of a loan
func (h *Handler) BorrowBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// ExtendLoan extends the loan of a book
func (h *Handler) ExtendLoan(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// ReturnBook returns the book
func (h *Handler) ReturnBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
