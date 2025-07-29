package handler

import (
	"net/http"
	"time"

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

// PingResponse represents the ping response
type PingResponse struct {
	Message   string `json:"message" example:"pong"`
	Timestamp string `json:"timestamp" example:"2023-07-29T17:30:00Z"`
	Status    string `json:"status" example:"ok"`
}
