package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/user/dto"
	userModel "github.com/nclsgg/despensa-digital/backend/internal/modules/user/model"
	"github.com/nclsgg/despensa-digital/backend/pkg/response"
)

type userHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) domain.UserHandler {
	return &userHandler{service}
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.UserResponseWrapper
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /user/{id} [get]
func (h *userHandler) GetUser(c *gin.Context) {
	rawID, ok := c.Get("userID")
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	id, ok := rawID.(uuid.UUID)
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserById(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		default:
			response.InternalError(c, "Failed to retrieve user")
		}
		return
	}

	response.OK(c, toUserResponse(user))
}

// GetCurrentUser returns the currently authenticated user
// @Summary Get current authenticated user
// @Tags User
// @Produce json
// @Success 200 {object} response.UserResponseWrapper
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /user/me [get]
func (h *userHandler) GetCurrentUser(c *gin.Context) {
	rawID, ok := c.Get("userID")
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	id, ok := rawID.(uuid.UUID)
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserById(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		default:
			response.InternalError(c, "Failed to retrieve user")
		}
		return
	}

	response.OK(c, toUserResponse(user))
}

// GetAllUsers returns all users (admin only)
// @Summary List all users (admin only)
// @Tags User
// @Produce json
// @Success 200 {array} response.UserListResponseWrapper
// @Failure 500 {object} response.APIResponse
// @Router /user/all [get]
func (h *userHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to retrieve users")
		return
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		user := users[i]
		responses = append(responses, toUserResponse(&user))
	}

	response.OK(c, responses)
}

// CompleteProfile completes user profile with name information
// @Summary Complete user profile after OAuth registration
// @Description Updates user's first name and last name
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CompleteProfileRequest true "Profile completion data"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /user/complete-profile [put]
func (h *userHandler) CompleteProfile(c *gin.Context) {
	var req dto.CompleteProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data")
		return
	}

	rawID, ok := c.Get("userID")
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	id, ok := rawID.(uuid.UUID)
	if !ok {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	err := h.service.CompleteProfile(c.Request.Context(), id, req.FirstName, req.LastName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		default:
			response.InternalError(c, "Failed to complete profile")
		}
		return
	}

	response.OK(c, gin.H{"message": "Profile completed successfully"})
}

func toUserResponse(user *userModel.User) dto.UserResponse {
	return dto.UserResponse{
		ID:               user.ID.String(),
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Role:             user.Role,
		ProfileCompleted: user.ProfileCompleted,
	}
}
