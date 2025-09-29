package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	profileService domain.ProfileService
}

func NewProfileHandler(profileService domain.ProfileService) (result0 *ProfileHandler) {
	__logParams := map[string]any{"profileService": profileService}
	__logStart := time.

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
		Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "NewProfileHandler"), zap.Any("result", result0), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "NewProfileHandler"), zap.Any("params", __logParams))
	result0 = &ProfileHandler{
		profileService: profileService,
	}
	return
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProfileHandler.CreateProfile"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProfileHandler.CreateProfile"), zap.Any("params", __logParams))
	var input dto.CreateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProfileHandler.CreateProfile"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.CreateProfile(c.Request.Context(), userID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProfileHandler.CreateProfile"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrProfileAlreadyExists):
			response.Fail(c, http.StatusConflict, "PROFILE_EXISTS", "Profile already exists")
		default:
			response.InternalError(c, "Failed to create profile")
		}
		return
	}

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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProfileHandler.GetProfile"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProfileHandler.GetProfile"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.GetProfileByUserID(c.Request.Context(), userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProfileHandler.GetProfile"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrProfileNotFound):
			response.Fail(c, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
		default:
			response.InternalError(c, "Failed to fetch profile")
		}
		return
	}

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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProfileHandler.UpdateProfile"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProfileHandler.UpdateProfile"), zap.Any("params", __logParams))
	var input dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProfileHandler.UpdateProfile"), zap.Error(err), zap.Any("params", __logParams))
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProfileHandler.UpdateProfile"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrProfileNotFound):
			response.Fail(c, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
		default:
			response.InternalError(c, "Failed to update profile")
		}
		return
	}

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
	__logParams := map[string]any{"h": h, "c": c}
	__logStart := time.Now()
	defer func() {
		zap.L().Info("function.exit", zap.String("func", "*ProfileHandler.DeleteProfile"), zap.Any("result", nil), zap.Duration("duration", time.Since(__logStart)))
	}()
	zap.L().Info("function.entry", zap.String("func", "*ProfileHandler.DeleteProfile"), zap.Any("params", __logParams))
	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err := h.profileService.DeleteProfile(c.Request.Context(), userID)
	if err != nil {
		zap.L().Error("function.error", zap.String("func", "*ProfileHandler.DeleteProfile"), zap.Error(err), zap.Any("params", __logParams))
		switch {
		case errors.Is(err, domain.ErrProfileNotFound):
			response.Fail(c, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
		default:
			response.InternalError(c, "Failed to delete profile")
		}
		return
	}

	response.OK(c, gin.H{"message": "Profile deleted successfully"})
}
