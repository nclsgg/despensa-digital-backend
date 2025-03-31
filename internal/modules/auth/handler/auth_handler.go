package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/dispensa-digital/backend/internal/modules/auth/model"
)

type authHandler struct {
	service domain.AuthService
}

func NewAuthHandler(service domain.AuthService) domain.AuthHandler {
	return &authHandler{service}
}

func (h *authHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.Register(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
