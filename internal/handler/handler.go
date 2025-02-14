package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, John")
}
