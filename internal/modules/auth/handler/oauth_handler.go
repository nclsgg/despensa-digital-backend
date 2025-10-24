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
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
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
// @Param mobile query string false "Mobile flag (true/false)"
// @Success 302 "Redirect to OAuth provider"
// @Router /auth/oauth/{provider} [get]
func (h *oauthHandler) OAuthLogin(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" {
		h.handleAuthError(c, domain.ErrProviderNotSupported)
		return
	}

	// Use OAuth state parameter to pass mobile flag reliably
	isMobile := c.Query("mobile") == "true"

	// Encode custom state with mobile flag
	stateData := map[string]interface{}{
		"mobile": isMobile,
		"nonce":  uuid.New().String(), // Add random nonce for security
	}
	stateJSON, _ := json.Marshal(stateData)
	customState := url.QueryEscape(string(stateJSON))

	// Store custom state in query for gothic to use
	query := c.Request.URL.Query()
	query.Set("state", customState)
	c.Request.URL.RawQuery = query.Encode()

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
	ctx := c.Request.Context()
	logger := appLogger.FromContext(ctx)

	provider := c.Param("provider")
	if provider != "google" {
		h.handleAuthError(c, domain.ErrProviderNotSupported)
		return
	}

	gothUser, err := h.ginGothicCompleteAuth(c)
	if err != nil {
		logger.Error("Failed to complete OAuth authentication",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "OAuthCallback"),
			zap.String(appLogger.FieldAction, "complete_auth"),
			zap.String("provider", provider),
			zap.Error(err),
		)
		response.InternalError(c, "Failed to complete auth")
		return
	}

	user, err := h.findOrCreateUser(ctx, gothUser)
	if err != nil {
		logger.Error("Failed to find or create user",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "OAuthCallback"),
			zap.String(appLogger.FieldAction, "find_or_create_user"),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(gothUser.Email)),
			zap.Error(err),
		)
		h.handleAuthError(c, err)
		return
	}

	accessToken, err := h.service.GenerateAccessToken(user)
	if err != nil {
		logger.Error("Failed to generate access token",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "OAuthCallback"),
			zap.String(appLogger.FieldAction, "generate_token"),
			zap.String(appLogger.FieldUserID, user.ID.String()),
			zap.Error(err),
		)
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

	// Extract mobile flag from OAuth state parameter
	isMobile := false
	if stateParam := c.Query("state"); stateParam != "" {
		stateJSON, err := url.QueryUnescape(stateParam)
		if err == nil {
			var stateData map[string]interface{}
			if err := json.Unmarshal([]byte(stateJSON), &stateData); err == nil {
				if mobile, ok := stateData["mobile"].(bool); ok {
					isMobile = mobile
				}
			}
		}
	}

	var redirectURL string

	if isMobile {
		redirectURL = fmt.Sprintf("despensadigital://auth/callback?data=%s", url.QueryEscape(string(authData)))
		logger.Info("OAuth callback redirecting to mobile app",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "OAuthCallback"),
			zap.String(appLogger.FieldAction, "redirect"),
			zap.String(appLogger.FieldUserID, user.ID.String()),
			zap.Bool("is_mobile", true),
		)
	} else {
		redirectURL = fmt.Sprintf("%s/auth/callback?data=%s", frontendURL, url.QueryEscape(string(authData)))
		logger.Info("OAuth callback redirecting to web frontend",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "OAuthCallback"),
			zap.String(appLogger.FieldAction, "redirect"),
			zap.String(appLogger.FieldUserID, user.ID.String()),
			zap.Bool("is_mobile", false),
		)
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *oauthHandler) findOrCreateUser(ctx context.Context, gothUser goth.User) (*model.User, error) {
	logger := appLogger.FromContext(ctx)

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

		logger.Info("Found existing OAuth user",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "findOrCreateUser"),
			zap.String(appLogger.FieldAction, "user_found"),
			zap.String(appLogger.FieldUserID, user.ID.String()),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(gothUser.Email)),
		)
		return user, nil
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		logger.Error("Error getting user by email",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "findOrCreateUser"),
			zap.String(appLogger.FieldAction, "get_user"),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(gothUser.Email)),
			zap.Error(err),
		)
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
		logger.Error("Failed to create OAuth user",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "findOrCreateUser"),
			zap.String(appLogger.FieldAction, "create_user"),
			zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(gothUser.Email)),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Created new OAuth user",
		zap.String(appLogger.FieldModule, "auth"),
		zap.String(appLogger.FieldFunction, "findOrCreateUser"),
		zap.String(appLogger.FieldAction, "user_created"),
		zap.String(appLogger.FieldUserID, newUser.ID.String()),
		zap.String(appLogger.FieldEmail, appLogger.SanitizeEmail(gothUser.Email)),
	)
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
	ctx := c.Request.Context()
	logger := appLogger.FromContext(ctx)

	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid profile update request",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "CompleteProfile"),
			zap.String(appLogger.FieldAction, "bind_json"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid request data")
		return
	}

	uid := userID.(uuid.UUID)
	if err := h.service.CompleteProfile(ctx, uid, req.FirstName, req.LastName); err != nil {
		logger.Error("Failed to complete profile",
			zap.String(appLogger.FieldModule, "auth"),
			zap.String(appLogger.FieldFunction, "CompleteProfile"),
			zap.String(appLogger.FieldAction, "complete_profile"),
			zap.String(appLogger.FieldUserID, uid.String()),
			zap.Error(err),
		)
		h.handleAuthError(c, err)
		return
	}

	logger.Info("Profile completed successfully",
		zap.String(appLogger.FieldModule, "auth"),
		zap.String(appLogger.FieldFunction, "CompleteProfile"),
		zap.String(appLogger.FieldAction, "complete_profile"),
		zap.String(appLogger.FieldUserID, uid.String()),
	)
	response.OK(c, gin.H{"message": "Profile updated successfully"})
}
