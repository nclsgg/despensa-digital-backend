package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

type authHandler struct {
	service domain.AuthService
}

func NewAuthHandler(service domain.AuthService) domain.AuthHandler {
	return &authHandler{service}
}

// Register registers a new user
// @Summary Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User data"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload")
		return
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	accessToken, refreshToken, err := h.service.Register(c.Request.Context(), user)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	authResp := dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.OK(c, authResp)
}

// Login authenticates a user and returns an access token
// @Summary User login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User credentials"
// @Success 200 {object} response.LoginSuccessResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload")
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.InternalError(c, "Failed to login")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "despensa-digital-backend-production.up.railway.app",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	authResp := dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.OK(c, authResp)
}

// Logout terminates the user session
// @Summary User logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	if err := h.service.Logout(c.Request.Context(), refreshToken); err != nil {
		response.InternalError(c, "Failed to logout")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "despensa-digital-backend-production.up.railway.app",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	response.OK(c, gin.H{"message": "Logout successful"})
}

// RefreshToken generates a new access token from the refresh token
// @Summary Refresh access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.LoginSuccessResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.BadRequest(c, "Invalid refresh token")
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.InternalError(c, "Failed to refresh token")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "despensa-digital-backend-production.up.railway.app",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	resp := dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	response.OK(c, resp)
}
