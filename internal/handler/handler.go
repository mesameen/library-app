package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/test/library-app/internal/logger"
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
		logger.Errorf("title is mandatory")
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is mandatory"})
		return
	}
	det, err := h.repo.GetBookDetails(c, title)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.Errorf("fetching title %s failed. Error: %v", title, err)
		// rest of all errors falls under this category
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, det)
}

// BorrowBook borrows a book from store (loan period: 4 weeks) and returns the details of a loan
func (h *Handler) BorrowBook(c *gin.Context) {
	var borrowReq model.LoanRequest
	err := c.BindJSON(&borrowReq)
	if err != nil {
		logger.Errorf("Failed to unamrshal the request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if borrowReq.NameOfBorrower == "" || borrowReq.Title == "" {
		logger.Errorf("NameOfBorrower & Title are mandatory to borrow a a book.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "NameOfBorrower or Title are missed in the request"})
		return
	}
	loanDetails := &model.LoanDetails{
		NameOfBorrower: borrowReq.NameOfBorrower,
		Title:          borrowReq.Title,
		LoanDate:       time.Now().Unix(),
		ReturnDate:     time.Now().Add(24 * 7 * time.Hour).Unix(), // 7 days return period
	}
	_, err = h.repo.AddLoan(c, loanDetails)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loanDetails)
}

// ExtendLoan extends the loan of a book
func (h *Handler) ExtendLoan(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Errorf("invalid id %s to update loan", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is mandatory"})
		return
	}
	// extenidng loan
	loan, err := h.repo.ExtendLoan(c, idInt)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// rest of all errors falls under this category
		logger.Errorf("fetching loan %d failed. Error: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"return_date": loan.ReturnDate, "message": "loan got extended to 3 weeks"})
}

// ReturnBook returns the book
func (h *Handler) ReturnBook(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Errorf("invalid id %s to update loan", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is mandatory"})
		return
	}
	err = h.repo.ReturnBook(c, idInt)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// rest of all errors falls under this category
		logger.Errorf("fetching loan %d failed. Error: %v", idInt, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "book returned"})
}
