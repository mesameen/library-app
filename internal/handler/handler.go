package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
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
			// unwrapping to send actual error
			err = errors.Unwrap(err)
			logger.Errorf("%s", err.Error())
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
	det, err := h.repo.GetBookDetails(c, borrowReq.Title)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			// unwrapping to send actual error
			uwErr := errors.Unwrap(err)
			logger.Errorf("%s", uwErr.Error())
			c.JSON(http.StatusNotFound, gin.H{"error": uwErr})
			return
		}
		// rest of all errors falls under this category
		logger.Errorf("fetching title %s failed. Error: %v", borrowReq.Title, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// if available copies are zero returning the error
	if det.AvailableCopies == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("not enough copies of requested title %v", borrowReq.Title)})
		return
	}
	loanDetails := &model.LoanDetails{
		ID:             GetUniqueIncrementedID(),
		NameOfBorrower: borrowReq.NameOfBorrower,
		Title:          borrowReq.Title,
		LoanDate:       time.Now().Unix(),
		ReturnDate:     time.Now().Add(24 * 7 * time.Hour).Unix(), // 7 days return period
	}
	_, err = h.repo.AddLoan(c, loanDetails)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("not enough copies of requested title %v", borrowReq.Title)})
		return
	}
	c.JSON(http.StatusCreated, loanDetails)
}

var uniqueID int32

func GetUniqueIncrementedID() int {
	// incrementing uniqueID
	atomic.AddInt32(&uniqueID, 1)
	return int(uniqueID)
}

// ExtendLoan extends the loan of a book
func (h *Handler) ExtendLoan(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// ReturnBook returns the book
func (h *Handler) ReturnBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
