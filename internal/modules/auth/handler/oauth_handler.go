package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/dto"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

type oauthHandler struct {
	service domain.AuthService
	cfg     *config.Config
}

func NewOAuthHandler(service domain.AuthService, cfg *config.Config) *oauthHandler {
	return &oauthHandler{service: service, cfg: cfg}
}

// InitOAuth initializes OAuth providers
func (h *oauthHandler) InitOAuth() {
	key := h.cfg.SessionSecret
	if len(key) < 32 {
		key = "fallback-session-secret-for-development-only-32-chars"
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(3600)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store

	goth.UseProviders(
		google.New(
			h.cfg.GoogleClientID,
			h.cfg.GoogleClientSecret,
			h.cfg.GoogleCallbackURL,
		),
	)

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		parts := strings.Split(req.URL.Path, "/")
		for i, part := range parts {
			if part == "oauth" && i+1 < len(parts) {
				return parts[i+1], nil
			}
		}
		return "google", nil
	}
}

// OAuthLogin initiates OAuth login
// @Summary Initiate OAuth login
// @Tags OAuth
// @Param provider path string true "OAuth Provider (google)"
// @Success 302 "Redirect to OAuth provider"
// @Router /auth/oauth/{provider} [get]
func (h *oauthHandler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" {
		h.handleAuthError(c, domain.ErrProviderNotSupported)
		return
	}

	h.ginGothicBeginAuth(c)
}

// OAuthCallback handles OAuth callback
// @Summary OAuth callback
// @Tags OAuth
// @Produce json
// @Param provider path string true "OAuth Provider (google)"
// @Param code query string true "OAuth code"
// @Param state query string true "OAuth state"
// @Success 200 {object} response.LoginSuccessResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/oauth/{provider}/callback [get]
func (h *oauthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" {
		h.handleAuthError(c, domain.ErrProviderNotSupported)
		return
	}

	gothUser, err := h.ginGothicCompleteAuth(c)
	if err != nil {
		response.InternalError(c, "Failed to complete auth")
		return
	}

	user, err := h.findOrCreateUser(c.Request.Context(), gothUser)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	accessToken, err := h.service.GenerateAccessToken(user)
	if err != nil {
		response.InternalError(c, "Failed to generate access token")
		return
	}

	authResp := dto.AuthResponse{
		AccessToken: accessToken,
		User: dto.UserDTO{
			ID:               user.ID.String(),
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Email:            user.Email,
			Role:             user.Role,
			ProfileCompleted: user.ProfileCompleted,
			IsActive:         true,
			CreatedAt:        user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        user.UpdatedAt.Format(time.RFC3339),
		},
	}

	frontendURL := h.cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	authData, _ := json.Marshal(authResp)
	redirectURL := fmt.Sprintf("%s/auth/callback?data=%s", frontendURL, url.QueryEscape(string(authData)))

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *oauthHandler) findOrCreateUser(ctx context.Context, gothUser goth.User) (*model.User, error) {
	user, err := h.service.GetUserByEmail(ctx, gothUser.Email)
	if err == nil {
		if gothUser.Name != "" && !user.ProfileCompleted {
			parts := strings.SplitN(gothUser.Name, " ", 2)
			if len(parts) > 0 {
				user.FirstName = parts[0]
			}
			if len(parts) > 1 {
				user.LastName = parts[1]
				user.ProfileCompleted = true
			}
		}
		return user, nil
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	newUser := &model.User{
		Email:            gothUser.Email,
		Role:             "user",
		ProfileCompleted: false,
	}

	if gothUser.Name != "" {
		parts := strings.SplitN(gothUser.Name, " ", 2)
		if len(parts) > 0 {
			newUser.FirstName = parts[0]
		}
		if len(parts) > 1 {
			newUser.LastName = parts[1]
			newUser.ProfileCompleted = true
		}
	}

	if err := h.service.CreateUserOAuth(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (h *oauthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrProviderNotSupported):
		response.BadRequest(c, "Provider not supported")
	case errors.Is(err, domain.ErrEmailAlreadyRegistered):
		response.Fail(c, http.StatusConflict, "CONFLICT", "Email already registered")
	case errors.Is(err, domain.ErrUserNotFound):
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", "User not found")
	case errors.Is(err, domain.ErrProfileUpdateFailed):
		response.InternalError(c, "Failed to update profile")
	default:
		response.InternalError(c, "Authentication error")
	}
}

func (h *oauthHandler) ginGothicBeginAuth(c *gin.Context) {
	writer := &ginResponseWriter{c.Writer, c}
	gothic.BeginAuthHandler(writer, c.Request)
}

func (h *oauthHandler) ginGothicCompleteAuth(c *gin.Context) (goth.User, error) {
	writer := &ginResponseWriter{c.Writer, c}
	return gothic.CompleteUserAuth(writer, c.Request)
}

type ginResponseWriter struct {
	gin.ResponseWriter
	ctx *gin.Context
}

func (w *ginResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ginResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func (w *ginResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// CompleteProfile completes user profile with first and last name
// @Summary Complete user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /auth/complete-profile [patch]
func (h *oauthHandler) CompleteProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	if err := h.service.CompleteProfile(c.Request.Context(), userID.(uuid.UUID), req.FirstName, req.LastName); err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "Profile updated successfully"})
}
