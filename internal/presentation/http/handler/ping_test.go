package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler_Ping(t *testing.T) {
	// Create a new Echo instance for testing
	e := echo.New()

	// Create a new ping handler
	handler := NewPingHandler()

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test the ping handler
	err := handler.Ping(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse response
	var response PingResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response fields
	assert.Equal(t, "pong", response.Message)
	assert.Equal(t, "ok", response.Status)

	// Verify timestamp is valid
	_, err = time.Parse(time.RFC3339, response.Timestamp)
	assert.NoError(t, err)
}

func TestNewPingHandler(t *testing.T) {
	// Test handler creation
	handler := NewPingHandler()

	// Verify handler is not nil
	assert.NotNil(t, handler)
	assert.IsType(t, &PingHandler{}, handler)
}

func TestPingResponse_Structure(t *testing.T) {
	// Test PingResponse struct
	response := PingResponse{
		Message:   "pong",
		Timestamp: "2023-07-29T17:30:00Z",
		Status:    "ok",
	}

	// Verify fields
	assert.Equal(t, "pong", response.Message)
	assert.Equal(t, "2023-07-29T17:30:00Z", response.Timestamp)
	assert.Equal(t, "ok", response.Status)

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// Verify JSON structure
	var unmarshaled PingResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, response, unmarshaled)
}
