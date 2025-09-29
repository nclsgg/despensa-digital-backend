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
	"go.uber.org/zap"
)

type oauthHandler struct {
	service domain.AuthService
	cfg     *config.Config
}

func NewOAuthHandler(service domain.AuthService, cfg *config.Config) (result0 *oauthHandler) {
	__logParams := map[string]any{"service": service, "cfg": cfg}
	__logStart :=

		// InitOAuth initializes OAuth providers
		time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewOAuthHandler"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewOAuthHandler"), zap.Any("params", __logParams))
	result0 = &oauthHandler{service: service, cfg: cfg}
	return
}

func (h *oauthHandler) InitOAuth() {
	__logParams := map[string]any{"h": h}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.InitOAuth"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.InitOAuth"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.OAuthLogin"), zap.Any("result", nil),

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
			zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.OAuthLogin"), zap.Any("params", __logParams))
	provider := c.Param("provider")
	if provider != "google" {
		h.handleAuthError(c, domain.ErrProviderNotSupported)
		return
	}

	h.ginGothicBeginAuth(c)
}

func (h *oauthHandler) OAuthCallback(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.OAuthCallback"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.OAuthCallback"), zap.Any("params", __logParams))
	provider := c.Param("provider")
	if provider != "google" {
		h.handleAuthError(c, domain.ErrProviderNotSupported)
		return
	}

	gothUser, err := h.ginGothicCompleteAuth(c)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*oauthHandler.OAuthCallback"), zap.Error(err), zap.Any("params", __logParams))
		response.InternalError(c, "Failed to complete auth")
		return
	}

	user, err := h.findOrCreateUser(c.Request.Context(), gothUser)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*oauthHandler.OAuthCallback"), zap.Error(err), zap.Any("params", __logParams))
		h.handleAuthError(c, err)
		return
	}

	accessToken, err := h.service.GenerateAccessToken(user)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*oauthHandler.OAuthCallback"), zap.Error(err), zap.Any("params", __logParams))
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

func (h *oauthHandler) findOrCreateUser(ctx context.Context, gothUser goth.User) (result0 *model.User, result1 error) {
	__logParams := map[string]any{"h": h, "ctx": ctx, "gothUser": gothUser}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.findOrCreateUser"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.findOrCreateUser"), zap.Any("params", __logParams))
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
		result0 = user
		result1 = nil
		return
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		result0 = nil
		result1 = err
		return
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
		zap.L().Error("function.error", zap.String("func", "*oauthHandler.findOrCreateUser"), zap.Error(err), zap.Any("params", __logParams))
		result0 = nil
		result1 = err
		return
	}
	result0 = newUser
	result1 = nil
	return
}

func (h *oauthHandler) handleAuthError(c *gin.Context, err error) {
	__logParams := map[string]any{"h": h, "c": c, "err": err}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.handleAuthError"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.handleAuthError"), zap.Any("params", __logParams))
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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.ginGothicBeginAuth"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.ginGothicBeginAuth"), zap.Any("params", __logParams))
	writer := &ginResponseWriter{c.Writer, c}
	gothic.BeginAuthHandler(writer, c.Request)
}

func (h *oauthHandler) ginGothicCompleteAuth(c *gin.Context) (result0 goth.User, result1 error) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.ginGothicCompleteAuth"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.ginGothicCompleteAuth"), zap.Any("params", __logParams))
	writer := &ginResponseWriter{c.Writer, c}
	result0, result1 = gothic.CompleteUserAuth(writer, c.Request)
	return
}

type ginResponseWriter struct {
	gin.ResponseWriter
	ctx *gin.Context
}

func (w *ginResponseWriter) Header() (result0 http.Header) {
	__logParams := map[string]any{"w": w}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ginResponseWriter.Header"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String(

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
		"func", "*ginResponseWriter.Header"), zap.Any("params", __logParams))
	result0 = w.ResponseWriter.Header()
	return
}

func (w *ginResponseWriter) Write(data []byte) (result0 int, result1 error) {
	__logParams := map[string]any{"w": w, "data": data}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ginResponseWriter.Write"), zap.Any("result", map[string]any{"result0": result0, "result1": result1}), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ginResponseWriter.Write"), zap.Any("params", __logParams))
	result0, result1 = w.ResponseWriter.Write(data)
	return
}

func (w *ginResponseWriter) WriteHeader(statusCode int) {
	__logParams := map[string]any{"w": w, "statusCode": statusCode}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ginResponseWriter.WriteHeader"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ginResponseWriter.WriteHeader"), zap.Any("params", __logParams))
	w.ResponseWriter.WriteHeader(statusCode)
}

func (h *oauthHandler) CompleteProfile(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*oauthHandler.CompleteProfile"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*oauthHandler.CompleteProfile"), zap.Any("params", __logParams))
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("function.error", zap.String("func", "*oauthHandler.CompleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid request data")
		return
	}

	if err := h.service.CompleteProfile(c.Request.Context(), userID.(uuid.UUID), req.FirstName, req.LastName); err != nil {
		zap.L().Error("function.error", zap.String("func", "*oauthHandler.CompleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "Profile updated successfully"})
}
