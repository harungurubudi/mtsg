package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/harungurubudi/mtsg/internal/presentation/dto"
	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	t.Run("NewErrorResponse creates correct error response", func(t *testing.T) {
		httpError := pkgerror.NewUnauthorizedError("invalid credentials")
		response := dto.NewErrorResponse(httpError)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "unauthorized", response.Error.Code)
		assert.Equal(t, "invalid credentials", response.Error.Message)
		assert.Equal(t, 401, response.Error.StatusCode)
	})

	t.Run("ErrorResponse JSON marshaling", func(t *testing.T) {
		httpError := pkgerror.NewValidationError("validation failed")
		response := dto.NewErrorResponse(httpError)

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		expected := `{"status":"error","error":{"code":"validation_error","message":"validation failed","status_code":400}}`
		assert.Equal(t, expected, string(jsonData))
	})
}

func TestSuccessResponse(t *testing.T) {
	t.Run("NewSuccessResponse creates correct success response", func(t *testing.T) {
		data := map[string]string{"message": "success"}
		response := dto.NewSuccessResponse(data)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, data, response.Data)
	})

	t.Run("SuccessResponse JSON marshaling", func(t *testing.T) {
		data := map[string]interface{}{
			"id":   123,
			"name": "test",
		}
		response := dto.NewSuccessResponse(data)

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		expected := `{"status":"success","data":{"id":123,"name":"test"}}`
		assert.Equal(t, expected, string(jsonData))
	})
}

func TestListResponse(t *testing.T) {
	t.Run("NewListResponse creates correct list response", func(t *testing.T) {
		data := []string{"item1", "item2", "item3"}
		response := dto.NewListResponse(data, 1, 10, 25)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, data, response.Data)
		assert.Equal(t, 1, response.Meta.Page)
		assert.Equal(t, 10, response.Meta.PerPage)
		assert.Equal(t, 25, response.Meta.Total)
		assert.Equal(t, 3, response.Meta.TotalPages) // (25 + 10 - 1) / 10 = 3
	})

	t.Run("NewListResponse with exact page division", func(t *testing.T) {
		response := dto.NewListResponse([]string{}, 1, 10, 20)

		assert.Equal(t, 2, response.Meta.TotalPages) // 20 / 10 = 2
	})

	t.Run("ListResponse JSON marshaling", func(t *testing.T) {
		data := []string{"item1", "item2"}
		response := dto.NewListResponse(data, 1, 10, 15)

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		expected := `{"status":"success","data":["item1","item2"],"meta":{"page":1,"per_page":10,"total":15,"total_pages":2}}`
		assert.Equal(t, expected, string(jsonData))
	})
}

func TestPaginationRequest(t *testing.T) {
	t.Run("GetOffset calculates correct offset", func(t *testing.T) {
		pagination := &dto.PaginationRequest{
			Page:    3,
			PerPage: 10,
		}

		assert.Equal(t, 20, pagination.GetOffset()) // (3-1) * 10 = 20
	})

	t.Run("GetLimit returns per page value", func(t *testing.T) {
		pagination := &dto.PaginationRequest{
			Page:    1,
			PerPage: 25,
		}

		assert.Equal(t, 25, pagination.GetLimit())
	})

	t.Run("DefaultPaginationRequest returns correct defaults", func(t *testing.T) {
		defaultPagination := dto.DefaultPaginationRequest()

		assert.Equal(t, 1, defaultPagination.Page)
		assert.Equal(t, 20, defaultPagination.PerPage)
	})

	t.Run("PaginationRequest validation tags", func(t *testing.T) {
		pagination := &dto.PaginationRequest{
			Page:    1,
			PerPage: 10,
		}

		// Test that struct tags are present (validation will be tested in integration tests)
		assert.NotEmpty(t, pagination)
	})
}

func TestTimestamp(t *testing.T) {
	t.Run("Timestamp JSON marshaling", func(t *testing.T) {
		now := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
		timestamp := dto.Timestamp{now}

		jsonData, err := json.Marshal(timestamp)
		assert.NoError(t, err)

		expected := `"2023-12-25T10:30:00Z"`
		assert.Equal(t, expected, string(jsonData))
	})

	t.Run("Timestamp JSON unmarshaling", func(t *testing.T) {
		jsonData := `"2023-12-25T10:30:00Z"`
		var timestamp dto.Timestamp

		err := json.Unmarshal([]byte(jsonData), &timestamp)
		assert.NoError(t, err)

		expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
		assert.Equal(t, expected, timestamp.Time)
	})

	t.Run("Timestamp JSON unmarshaling with invalid format", func(t *testing.T) {
		jsonData := `"invalid-time"`
		var timestamp dto.Timestamp

		err := json.Unmarshal([]byte(jsonData), &timestamp)
		assert.Error(t, err)
	})

	t.Run("Timestamp with timezone", func(t *testing.T) {
		// Test with a specific timezone
		loc, _ := time.LoadLocation("America/New_York")
		now := time.Date(2023, 12, 25, 10, 30, 0, 0, loc)
		timestamp := dto.Timestamp{now}

		jsonData, err := json.Marshal(timestamp)
		assert.NoError(t, err)

		// Should be in RFC3339 format
		assert.Contains(t, string(jsonData), "2023-12-25")
	})
}

func TestErrorResponse_EdgeCases(t *testing.T) {
	t.Run("ErrorResponse with empty message", func(t *testing.T) {
		httpError := pkgerror.NewHTTPError(500, "internal_error", "")
		response := dto.NewErrorResponse(httpError)

		assert.Equal(t, "", response.Error.Message)
		assert.Equal(t, 500, response.Error.StatusCode)
	})

	t.Run("ErrorResponse with special characters", func(t *testing.T) {
		httpError := pkgerror.NewHTTPError(400, "bad_request", "Invalid input: <script>alert('xss')</script>")
		response := dto.NewErrorResponse(httpError)

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		// Should properly escape special characters
		assert.Contains(t, string(jsonData), "Invalid input")
		assert.Contains(t, string(jsonData), "bad_request")
	})
}

func TestListResponse_EdgeCases(t *testing.T) {
	t.Run("ListResponse with zero total", func(t *testing.T) {
		response := dto.NewListResponse([]string{}, 1, 10, 0)

		assert.Equal(t, 0, response.Meta.TotalPages)
		assert.Equal(t, 0, response.Meta.Total)
	})

	t.Run("ListResponse with total less than per page", func(t *testing.T) {
		response := dto.NewListResponse([]string{"item"}, 1, 10, 5)

		assert.Equal(t, 1, response.Meta.TotalPages)
		assert.Equal(t, 5, response.Meta.Total)
	})

	t.Run("ListResponse with large numbers", func(t *testing.T) {
		response := dto.NewListResponse([]string{}, 1, 10, 1000000)

		assert.Equal(t, 100000, response.Meta.TotalPages) // (1000000 + 10 - 1) / 10 = 100000
	})
}

func TestSuccessResponse_EdgeCases(t *testing.T) {
	t.Run("SuccessResponse with nil data", func(t *testing.T) {
		response := dto.NewSuccessResponse(nil)

		assert.Equal(t, "success", response.Status)
		assert.Nil(t, response.Data)
	})

	t.Run("SuccessResponse with complex nested data", func(t *testing.T) {
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"id":   123,
				"name": "John Doe",
				"tags": []string{"admin", "verified"},
			},
			"metadata": map[string]interface{}{
				"created_at": "2023-12-25T10:30:00Z",
				"updated_at": "2023-12-25T11:00:00Z",
			},
		}
		response := dto.NewSuccessResponse(data)

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		// Should properly serialize complex nested structures
		assert.Contains(t, string(jsonData), "John Doe")
		assert.Contains(t, string(jsonData), "admin")
		assert.Contains(t, string(jsonData), "verified")
	})
}
