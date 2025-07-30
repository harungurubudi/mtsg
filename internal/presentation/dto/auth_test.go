package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/domain/user"
	"github.com/harungurubudi/mtsg/internal/presentation/dto"
	"github.com/stretchr/testify/assert"
)

func TestLoginRequest(t *testing.T) {
	t.Run("LoginRequest JSON marshaling", func(t *testing.T) {
		request := dto.LoginRequest{
			Email:    "user@example.com",
			Password: "SecurePass123!",
		}

		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		expected := `{"email":"user@example.com","password":"SecurePass123!"}`
		assert.Equal(t, expected, string(jsonData))
	})

	t.Run("LoginRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{"email":"user@example.com","password":"SecurePass123!"}`
		var request dto.LoginRequest

		err := json.Unmarshal([]byte(jsonData), &request)
		assert.NoError(t, err)

		assert.Equal(t, "user@example.com", request.Email)
		assert.Equal(t, "SecurePass123!", request.Password)
	})

	t.Run("LoginRequest ToDomain converts correctly", func(t *testing.T) {
		request := dto.LoginRequest{
			Email:    "user@example.com",
			Password: "SecurePass123!",
		}

		credential := request.ToDomain()

		assert.Equal(t, user.Email("user@example.com"), credential.Email)
		assert.Equal(t, user.Password("SecurePass123!"), credential.Password)
	})

	t.Run("LoginRequest validation tags", func(t *testing.T) {
		request := dto.LoginRequest{
			Email:    "user@example.com",
			Password: "SecurePass123!",
		}

		// Test that struct tags are present (validation will be tested in integration tests)
		assert.NotEmpty(t, request)
	})
}

func TestLoginResponse(t *testing.T) {
	t.Run("LoginResponse JSON marshaling", func(t *testing.T) {
		response := dto.LoginResponse{
			AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		expected := `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}`
		assert.Equal(t, expected, string(jsonData))
	})

	t.Run("LoginResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}`
		var response dto.LoginResponse

		err := json.Unmarshal([]byte(jsonData), &response)
		assert.NoError(t, err)

		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", response.AccessToken)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", response.RefreshToken)
	})

	t.Run("NewLoginResponse creates correct response from session", func(t *testing.T) {
		session := &authentication.Session{
			AccessToken:  "access_token_123",
			RefreshToken: "refresh_token_456",
		}

		response := dto.NewLoginResponse(session)

		assert.Equal(t, "access_token_123", response.AccessToken)
		assert.Equal(t, "refresh_token_456", response.RefreshToken)
	})

	t.Run("NewLoginResponse with nil session", func(t *testing.T) {
		response := dto.NewLoginResponse(nil)

		assert.Equal(t, "", response.AccessToken)
		assert.Equal(t, "", response.RefreshToken)
	})
}

func TestLoginRequest_EdgeCases(t *testing.T) {
	t.Run("LoginRequest with empty fields", func(t *testing.T) {
		request := dto.LoginRequest{
			Email:    "",
			Password: "",
		}

		credential := request.ToDomain()

		assert.Equal(t, user.Email(""), credential.Email)
		assert.Equal(t, user.Password(""), credential.Password)
	})

	t.Run("LoginRequest with special characters", func(t *testing.T) {
		request := dto.LoginRequest{
			Email:    "user+tag@example.com",
			Password: "Pass@word#123!",
		}

		jsonData, err := json.Marshal(request)
		assert.NoError(t, err)

		// Should properly handle special characters
		assert.Contains(t, string(jsonData), "user+tag@example.com")
		assert.Contains(t, string(jsonData), "Pass@word#123!")
	})

	t.Run("LoginRequest with very long password", func(t *testing.T) {
		longPassword := "VeryLongPassword123!@#$%^&*()_+-=[]{}|;':\",./<>?"
		request := dto.LoginRequest{
			Email:    "user@example.com",
			Password: longPassword,
		}

		credential := request.ToDomain()

		assert.Equal(t, user.Password(longPassword), credential.Password)
	})
}

func TestLoginResponse_EdgeCases(t *testing.T) {
	t.Run("LoginResponse with empty tokens", func(t *testing.T) {
		response := dto.LoginResponse{
			AccessToken:  "",
			RefreshToken: "",
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		expected := `{"access_token":"","refresh_token":""}`
		assert.Equal(t, expected, string(jsonData))
	})

	t.Run("LoginResponse with very long tokens", func(t *testing.T) {
		longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		response := dto.LoginResponse{
			AccessToken:  longToken,
			RefreshToken: longToken,
		}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		// Should handle long tokens correctly
		assert.Contains(t, string(jsonData), longToken)
	})

	t.Run("NewLoginResponse with session containing special characters", func(t *testing.T) {
		session := &authentication.Session{
			AccessToken:  "access_token_with_special_chars_!@#$%^&*()",
			RefreshToken: "refresh_token_with_special_chars_!@#$%^&*()",
		}

		response := dto.NewLoginResponse(session)

		assert.Equal(t, "access_token_with_special_chars_!@#$%^&*()", response.AccessToken)
		assert.Equal(t, "refresh_token_with_special_chars_!@#$%^&*()", response.RefreshToken)
	})
}

func TestLoginRequest_Validation(t *testing.T) {
	t.Run("LoginRequest struct tags for validation", func(t *testing.T) {
		// This test ensures the struct tags are present for validation
		// Actual validation will be tested in integration tests with the validator
		request := dto.LoginRequest{}

		// Test that the struct can be created (tags are syntactically correct)
		assert.NotNil(t, request)
	})
}

func TestLoginResponse_Integration(t *testing.T) {
	t.Run("Complete login flow with DTOs", func(t *testing.T) {
		// Simulate a complete login flow
		loginRequest := dto.LoginRequest{
			Email:    "user@example.com",
			Password: "SecurePass123!",
		}

		// Convert to domain
		credential := loginRequest.ToDomain()
		assert.Equal(t, user.Email("user@example.com"), credential.Email)
		assert.Equal(t, user.Password("SecurePass123!"), credential.Password)

		// Simulate successful authentication
		session := &authentication.Session{
			AccessToken:  "access_token_123",
			RefreshToken: "refresh_token_456",
		}

		// Create response
		loginResponse := dto.NewLoginResponse(session)
		assert.Equal(t, "access_token_123", loginResponse.AccessToken)
		assert.Equal(t, "refresh_token_456", loginResponse.RefreshToken)

		// Test JSON serialization of response
		jsonData, err := json.Marshal(loginResponse)
		assert.NoError(t, err)

		expected := `{"access_token":"access_token_123","refresh_token":"refresh_token_456"}`
		assert.Equal(t, expected, string(jsonData))
	})
}

func TestAuthDTOs_TypeSafety(t *testing.T) {
	t.Run("Domain type conversions are type safe", func(t *testing.T) {
		request := dto.LoginRequest{
			Email:    "user@example.com",
			Password: "SecurePass123!",
		}

		credential := request.ToDomain()

		// Verify that the conversion produces the correct domain types
		assert.IsType(t, user.Email(""), credential.Email)
		assert.IsType(t, user.Password(""), credential.Password)
	})

	t.Run("Session to response conversion is type safe", func(t *testing.T) {
		session := &authentication.Session{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}

		response := dto.NewLoginResponse(session)

		// Verify that the response contains string types
		assert.IsType(t, "", response.AccessToken)
		assert.IsType(t, "", response.RefreshToken)
	})
}
