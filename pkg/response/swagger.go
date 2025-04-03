package response

// MessagePayload is used for generic success messages
type MessagePayload struct {
	Message string `json:"message"`
}

// MessageResponse is the generic wrapper for simple message responses
type MessageResponse struct {
	Success bool           `json:"success"`
	Data    MessagePayload `json:"data"`
	Error   *APIError      `json:"error,omitempty"`
}

// LoginSuccessResponse represents the structure of a successful login
type LoginSuccessResponse struct {
	Success bool         `json:"success"`
	Data    LoginPayload `json:"data"`
	Error   *APIError    `json:"error,omitempty"`
}

type LoginPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserResponseWrapper is used for GetUser or /me
type UserResponseWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *APIError   `json:"error,omitempty"`
}

// UserListResponseWrapper is used for GetAllUsers
type UserListResponseWrapper struct {
	Success bool          `json:"success"`
	Data    []interface{} `json:"data"`
	Error   *APIError     `json:"error,omitempty"`
}
