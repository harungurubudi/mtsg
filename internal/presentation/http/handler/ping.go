package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harungurubudi/mtsg/internal/presentation/http/middleware"
	"github.com/labstack/echo/v4"
)

// PingHandler handles ping-related endpoints
type PingHandler struct{}

// NewPingHandler creates a new PingHandler instance
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping godoc
// @Summary Health check endpoint
// @Description Simple ping endpoint to check if the server is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse "Server is running"
// @Router /ping [get]
func (h *PingHandler) Ping(c echo.Context) error {
	response := PingResponse{
		Message:   "pong",
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "ok",
	}

	return c.JSON(http.StatusOK, response)
}

// ProtectedPing godoc
// @Summary Protected health check endpoint
// @Description Protected ping endpoint that requires authentication
// @Tags Health
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProtectedPingResponse "Protected ping successful"
// @Failure 401 {object} error.HTTPError "Unauthorized - missing or invalid token"
// @Router /ping/protected [get]
func (h *PingHandler) ProtectedPing(c echo.Context) error {
	// Get user from context (set by auth middleware)
	user, exists := middleware.GetUserFromContext(c)
	if !exists || user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "user not authenticated",
		})
	}

	response := ProtectedPingResponse{
		Message:   "protected pong",
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "ok",
		UserID:    uuid.UUID(user.ID).String(),
		UserEmail: string(user.Email),
	}

	return c.JSON(http.StatusOK, response)
}

// PingResponse represents the ping response
type PingResponse struct {
	Message   string `json:"message" example:"pong"`
	Timestamp string `json:"timestamp" example:"2023-07-29T17:30:00Z"`
	Status    string `json:"status" example:"ok"`
}

// ProtectedPingResponse represents the protected ping response
type ProtectedPingResponse struct {
	Message   string `json:"message" example:"protected pong"`
	Timestamp string `json:"timestamp" example:"2023-07-29T17:30:00Z"`
	Status    string `json:"status" example:"ok"`
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserEmail string `json:"user_email" example:"user@example.com"`
}
