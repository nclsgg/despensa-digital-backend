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

func (h *authHandler) Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 60*60*24*7, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
