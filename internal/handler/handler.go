package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/test/library-app/internal/constants"
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

func (h *Handler) Live(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func (h *Handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "")
}

// GetAllBooks godoc
//
//	@Summary 		GetAllBooks fetches the book details
//	@Description 	GetAllBooks retrieves the detail and available copies of a book title
//	@Produce 		json
//	@Success 		200	{array}	model.BookDetails
//	@Failure 		404	{object}	model.CustomError
//	@Failure 		500	{object}	model.CustomError
//	@Router 		/book	[get]
//
// GetAllBooks retrieves all books in store
func (h *Handler) GetAllBooks(c *gin.Context) {
	det, err := h.repo.GetAllBookDetails(c)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			customError := &model.CustomError{
				Error: "zero books available in store",
				Code:  http.StatusNotFound,
			}
			c.JSON(http.StatusNotFound, customError)
			return
		}
		// rest of all errors falls under this category
		customError := &model.CustomError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, customError)
		return
	}
	c.JSON(http.StatusOK, det)
}

// GetAllLoans godoc
//
//	@Summary 		GetAllLoans fetches the all loan details
//	@Description 	GetAllLoans retrieves the detail of all loans
//	@Produce 		json
//	@Success 		200	{array}		model.LoanDetails
//	@Failure 		404	{object}	model.CustomError
//	@Failure 		500	{object}	model.CustomError
//	@Router 		/loan	[get]
//
// GetAllLoans retrieves all loans from store
func (h *Handler) GetAllLoans(c *gin.Context) {
	det, err := h.repo.GetAllLoans(c)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			customError := &model.CustomError{
				Error: "zero loans available in store",
				Code:  http.StatusNotFound,
			}
			c.JSON(http.StatusNotFound, customError)
			return
		}
		// rest of all errors falls under this category
		customError := &model.CustomError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, customError)
		return
	}
	if len(det) == 0 {
		customError := &model.CustomError{
			Error: "zero loans available in store",
			Code:  http.StatusNotFound,
		}
		c.JSON(http.StatusNotFound, customError)
		return
	}
	c.JSON(http.StatusOK, det)
}

// GetBook godoc
//
//	@Summary 		GetBook fetches the book details
//	@Description 	GetBook retrieves the detail and available copies of a book title
//	@Param			title	path	string	true	"Title of the book"
//	@Produce 		json
//	@Success 		200	{object}	model.BookDetails
//	@Failure 		404	{object}	model.CustomError
//	@Failure 		500	{object}	model.CustomError
//	@Router 		/book/{title}	[get]
//
// GetBook retrieves the detail and available copies of a book title
func (h *Handler) GetBook(c *gin.Context) {
	title := c.Param("title")
	if len(title) == 0 {
		logger.Errorf("title is mandatory")
		customError := &model.CustomError{
			Error: "title is mandatory",
			Code:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, customError)
		return
	}
	det, err := h.repo.GetBookDetails(c, title)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			customError := &model.CustomError{
				Error: err.Error(),
				Code:  http.StatusNotFound,
			}
			c.JSON(http.StatusNotFound, customError)
			return
		}
		logger.Errorf("fetching title %s failed. Error: %v", title, err)
		// rest of all errors falls under this category
		customError := &model.CustomError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, customError)
		return
	}
	c.JSON(http.StatusOK, det)
}

// LoanBook godoc
//
//	@Summary 		LoanBook borrows a book from store
//	@Description 	LoanBook borrows a book from store (loan period: 4 weeks) and returns the details of a loan
//	@Param			loanRequest	body	model.LoanRequest	true "Loan Request"
//	@Consume 		json	model.LoanRequest
//	@Produce 		json
//	@Success 		201	{object}	model.LoanDetails
//	@Failure 		404	{object}	model.CustomError
//	@Failure 		400	{object}	model.CustomError
//	@Failure 		500	{object}	model.CustomError
//	@Router 		/loan	[post]
//
// LoanBook borrows a book from store (loan period: 4 weeks) and returns the details of a loan
func (h *Handler) LoanBook(c *gin.Context) {
	var borrowReq model.LoanRequest
	err := c.BindJSON(&borrowReq)
	if err != nil {
		logger.Errorf("Failed to unamrshal the request body: %v", err)
		customError := &model.CustomError{
			Error: "internal server error",
			Code:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, customError)
		return
	}
	if borrowReq.NameOfBorrower == "" || borrowReq.Title == "" {
		logger.Errorf("NameOfBorrower & Title are mandatory to borrow a a book.")
		customError := &model.CustomError{
			Error: "NameOfBorrower or Title missed in the request",
			Code:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, customError)
		return
	}
	loanDetails := &model.LoanDetails{
		NameOfBorrower: borrowReq.NameOfBorrower,
		Title:          borrowReq.Title,
		LoanDate:       time.Now().Unix(),
		ReturnDate:     time.Now().Add(4 * 7 * 24 * time.Hour).Unix(), // 4 weeks return period
		Status:         constants.Active,
	}
	_, err = h.repo.AddLoan(c, loanDetails)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			customError := &model.CustomError{
				Error: err.Error(),
				Code:  http.StatusNotFound,
			}
			c.JSON(http.StatusNotFound, customError)
			return
		}
		customError := &model.CustomError{
			Error: "adding loan failed",
			Code:  http.StatusConflict,
		}
		c.JSON(http.StatusConflict, customError)
		return
	}
	c.JSON(http.StatusCreated, loanDetails)
}

// ExtendLoan godoc
//
//	@Summary 		ExtendLoan extends the loan of a book
//	@Description 	ExtendLoan extends the loan of a book
//	@Param			id	path	int	true	"Loan id"
//	@Consume 		json	model.LoanRequest
//	@Produce 		json
//	@Success 		202	{object}	model.LoanDetails
//	@Failure 		404	{object}	model.CustomError
//	@Failure 		400	{object}	model.CustomError
//	@Router 		/loan/extend/{id}	[post]
//
// ExtendLoan extends the loan of a book
func (h *Handler) ExtendLoan(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Errorf("invalid id %s to update loan", id)
		customError := &model.CustomError{
			Error: "id is mandatory",
			Code:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, customError)
		return
	}
	// extenidng loan
	loan, err := h.repo.ExtendLoan(c, idInt)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			customError := &model.CustomError{
				Error: err.Error(),
				Code:  http.StatusNotFound,
			}
			c.JSON(http.StatusNotFound, customError)
			return
		}
		// rest of all errors falls under this category
		logger.Errorf("fetching loan %d failed. Error: %v", idInt, err)
		customError := &model.CustomError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, customError)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"return_date": loan.ReturnDate, "message": "loan got extended to 3 weeks"})
}

// ReturnBook godoc
//
//	@Summary 		ReturnBook returns the book
//	@Description 	ReturnBook returns the book
//	@Param			id	path	int	true	"Loan id"
//	@Produce 		json
//	@Success 		202	{object}	model.LoanDetails
//	@Failure 		404	{object}	model.CustomError
//	@Failure 		400	{object}	model.CustomError
//	@Router 		/loan/return/{id}	[post]
//
// ReturnBook returns the book
func (h *Handler) ReturnBook(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Errorf("invalid id %s to update loan", id)
		customError := &model.CustomError{
			Error: "id is mandatory",
			Code:  http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, customError)
		return
	}
	err = h.repo.ReturnBook(c, idInt)
	if err != nil {
		// if notfound needs to return the specific error code and details
		if errors.Is(err, model.ErrNotFound) {
			customError := &model.CustomError{
				Error: err.Error(),
				Code:  http.StatusNotFound,
			}
			c.JSON(http.StatusNotFound, customError)
			return
		}
		// rest of all errors falls under this category
		logger.Errorf("fetching loan %d failed. Error: %v", idInt, err)
		customError := &model.CustomError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, customError)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "book returned"})
}
