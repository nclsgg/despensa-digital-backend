package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	// Configure session store for Gothic with proper settings
	key := h.cfg.SessionSecret
	if len(key) < 32 {
		key = "fallback-session-secret-for-development-only-32-chars"
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(3600) // 1 hour for OAuth flow
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store

	// Configure OAuth providers
	goth.UseProviders(
		google.New(
			h.cfg.GoogleClientID,
			h.cfg.GoogleClientSecret,
			h.cfg.GoogleCallbackURL,
		),
	)

	// Configure Gothic for Gin - extract provider from URL path
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		// Extract provider from URL path like /auth/oauth/google or /auth/oauth/google/callback
		parts := strings.Split(req.URL.Path, "/")
		for i, part := range parts {
			if part == "oauth" && i+1 < len(parts) {
				return parts[i+1], nil
			}
		}
		return "google", nil // fallback
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
		response.BadRequest(c, "Provider not supported")
		return
	}

	fmt.Printf("DEBUG: Starting OAuth login for provider: %s\n", provider)
	fmt.Printf("DEBUG: Request URL: %s\n", c.Request.URL.String())

	// Use Gin-compatible Gothic handler
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
		response.BadRequest(c, "Provider not supported")
		return
	}

	fmt.Printf("DEBUG: OAuth callback for provider: %s\n", provider)
	fmt.Printf("DEBUG: Request URL: %s\n", c.Request.URL.String())
	fmt.Printf("DEBUG: Request headers: %+v\n", c.Request.Header)

	// Complete OAuth authentication using Gin-compatible handler
	gothUser, err := h.ginGothicCompleteAuth(c)
	if err != nil {
		fmt.Printf("DEBUG: CompleteUserAuth error: %v\n", err)
		response.InternalError(c, fmt.Sprintf("Failed to complete auth: %v", err))
		return
	}

	fmt.Printf("DEBUG: Gothic user: %+v\n", gothUser)

	// Find or create user
	user, err := h.findOrCreateUser(gothUser)
	if err != nil {
		response.InternalError(c, "Failed to process user")
		return
	}

	// Generate access token for API access
	accessToken, err := h.service.GenerateAccessToken(user)
	if err != nil {
		response.InternalError(c, "Failed to generate access token")
		return
	}

	// Prepare response
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
			CreatedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	// Always redirect to frontend with token data
	frontendURL := h.cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Encode auth response as URL parameter
	authData, _ := json.Marshal(authResp)
	redirectURL := fmt.Sprintf("%s/auth/callback?data=%s", frontendURL, url.QueryEscape(string(authData)))
	
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// findOrCreateUser finds existing user or creates new one from OAuth data
func (h *oauthHandler) findOrCreateUser(gothUser goth.User) (*model.User, error) {
	// Try to find existing user by email
	existingUser, err := h.service.GetUserByEmail(gothUser.Email)
	if err == nil {
		// Update user info from OAuth if name is provided and profile is not completed
		if gothUser.Name != "" && !existingUser.ProfileCompleted {
			// Try to split the name into first and last name
			nameParts := strings.SplitN(gothUser.Name, " ", 2)
			if len(nameParts) > 0 {
				existingUser.FirstName = nameParts[0]
			}
			if len(nameParts) > 1 {
				existingUser.LastName = nameParts[1]
				existingUser.ProfileCompleted = true
			}
		}
		return existingUser, nil
	}

	// Create new user
	newUser := &model.User{
		Email: gothUser.Email,
		Role:  "user", // Default role
		ProfileCompleted: false, // Will need to complete profile
	}

	// Try to extract first and last name from OAuth if available
	if gothUser.Name != "" {
		nameParts := strings.SplitN(gothUser.Name, " ", 2)
		if len(nameParts) > 0 {
			newUser.FirstName = nameParts[0]
		}
		if len(nameParts) > 1 {
			newUser.LastName = nameParts[1]
			newUser.ProfileCompleted = true
		}
	}

	err = h.service.CreateUserOAuth(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// ginGothicBeginAuth adapts Gothic BeginAuthHandler for Gin
func (h *oauthHandler) ginGothicBeginAuth(c *gin.Context) {
	// Create a wrapper that implements http.ResponseWriter properly for Gin
	writer := &ginResponseWriter{c.Writer, c}
	
	// Call Gothic BeginAuthHandler with proper writer
	gothic.BeginAuthHandler(writer, c.Request)
}

// ginGothicCompleteAuth adapts Gothic CompleteUserAuth for Gin
func (h *oauthHandler) ginGothicCompleteAuth(c *gin.Context) (goth.User, error) {
	// Create a wrapper that implements http.ResponseWriter properly for Gin
	writer := &ginResponseWriter{c.Writer, c}
	
	// Call Gothic CompleteUserAuth with proper writer
	return gothic.CompleteUserAuth(writer, c.Request)
}

// ginResponseWriter wraps Gin's ResponseWriter to be compatible with http.ResponseWriter
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
	// Get user ID from context (set by authentication middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse request body
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	// Update profile
	err := h.service.CompleteProfile(c.Request.Context(), userID.(uuid.UUID), req.FirstName, req.LastName)
	if err != nil {
		response.InternalError(c, "Failed to update profile")
		return
	}

	response.OK(c, gin.H{"message": "Profile updated successfully"})
}
