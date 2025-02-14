package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/test/library-app/internal/store"
)

type Handler struct {
	repo store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{
		repo: s,
	}
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, John")
}

func (h *Handler) GetBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) BorrowBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) ExtendBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *Handler) ReturnBook(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
