package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/profile/dto"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
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
	var input dto.CreateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.CreateProfile(c.Request.Context(), userID, input)
	if err != nil {
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
	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.GetProfileByUserID(c.Request.Context(), userID)
	if err != nil {
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
	var input dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input: "+err.Error())
		return
	}

	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	profile, err := h.profileService.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
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
	rawID, _ := c.Get("userID")
	userID := rawID.(uuid.UUID)

	err := h.profileService.DeleteProfile(c.Request.Context(), userID)
	if err != nil {
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
