package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	appLogger "github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	profileService domain.ProfileService
}

func NewProfileHandler(profileService domain.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// CreateProfile godoc
// @Summary Create user profile
// @Description Create a new profile for the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Param profile body dto.CreateProfileDTO true "Profile data"
// @Success 201 {object} response.APIResponse{data=dto.ProfileResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile [post]
// @Security BearerAuth
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	var input dto.CreateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Warn("Invalid profile creation request",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "CreateProfile"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.CreateProfile(c.Request.Context(), userID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProfileAlreadyExists):
			logger.Warn("Attempt to create duplicate profile",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "CreateProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusConflict, "PROFILE_EXISTS", "Profile already exists")
		default:
			logger.Error("Failed to create profile",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "CreateProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to create profile")
		}
		return
	}

	logger.Info("Profile created via handler",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "CreateProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

	response.Success(c, http.StatusCreated, profile)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get profile for the authenticated user
// @Tags profile
// @Produce json
// @Success 200 {object} response.APIResponse{data=dto.ProfileResponseDTO}
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile [get]
// @Security BearerAuth
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.GetProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProfileNotFound):
			logger.Warn("Profile not found",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "GetProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
		default:
			logger.Error("Failed to fetch profile",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "GetProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to fetch profile")
		}
		return
	}

	logger.Info("Profile retrieved via handler",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "GetProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

	response.OK(c, profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update profile for the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Param profile body dto.UpdateProfileDTO true "Profile update data"
// @Success 200 {object} response.APIResponse{data=dto.ProfileResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile [put]
// @Security BearerAuth
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	var input dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Warn("Invalid profile update request",
			zap.String(appLogger.FieldModule, "profile"),
			zap.String(appLogger.FieldFunction, "UpdateProfile"),
			zap.Error(err),
		)
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProfileNotFound):
			logger.Warn("Profile not found for update",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "UpdateProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
		default:
			logger.Error("Failed to update profile",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "UpdateProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to update profile")
		}
		return
	}

	logger.Info("Profile updated via handler",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "UpdateProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

	response.OK(c, profile)
}

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Delete profile for the authenticated user
// @Tags profile
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /profile [delete]
// @Security BearerAuth
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	logger := appLogger.FromContext(c.Request.Context())

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err := h.profileService.DeleteProfile(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProfileNotFound):
			logger.Warn("Profile not found for deletion",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "DeleteProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
			)
			response.Fail(c, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
		default:
			logger.Error("Failed to delete profile",
				zap.String(appLogger.FieldModule, "profile"),
				zap.String(appLogger.FieldFunction, "DeleteProfile"),
				zap.String(appLogger.FieldUserID, userID.String()),
				zap.Error(err),
			)
			response.InternalError(c, "Failed to delete profile")
		}
		return
	}

	logger.Info("Profile deleted via handler",
		zap.String(appLogger.FieldModule, "profile"),
		zap.String(appLogger.FieldFunction, "DeleteProfile"),
		zap.String(appLogger.FieldUserID, userID.String()),
	)

	response.OK(c, gin.H{"message": "Profile deleted successfully"})
}
