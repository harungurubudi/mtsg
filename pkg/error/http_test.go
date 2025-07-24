package error_test

import (
	"testing"

	"github.com/harungurubudi/mtsg/pkg/error"
	"github.com/stretchr/testify/suite"
)

type HTTPErrorTestSuite struct {
	suite.Suite
}

func TestHTTPErrorTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPErrorTestSuite))
}

func (suite *HTTPErrorTestSuite) TestHTTPError_Error() {
	tests := []struct {
		name     string
		httpErr  *error.HTTPError
		expected string
	}{
		{
			name: "ShouldReturnMessage",
			httpErr: &error.HTTPError{
				StatusCode: 404,
				Code:       "not_found",
				Message:    "Resource not found",
			},
			expected: "Resource not found",
		},
		{
			name: "ShouldReturnEmptyMessage",
			httpErr: &error.HTTPError{
				StatusCode: 500,
				Code:       "internal_error",
				Message:    "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := tt.httpErr.Error()
			suite.Equal(tt.expected, result)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestHTTPError_GetStatusCode() {
	tests := []struct {
		name     string
		httpErr  *error.HTTPError
		expected int
	}{
		{
			name: "ShouldReturn404",
			httpErr: &error.HTTPError{
				StatusCode: 404,
				Code:       "not_found",
				Message:    "Not found",
			},
			expected: 404,
		},
		{
			name: "ShouldReturn500",
			httpErr: &error.HTTPError{
				StatusCode: 500,
				Code:       "internal_error",
				Message:    "Internal error",
			},
			expected: 500,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := tt.httpErr.GetStatusCode()
			suite.Equal(tt.expected, result)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestHTTPError_GetCode() {
	tests := []struct {
		name     string
		httpErr  *error.HTTPError
		expected string
	}{
		{
			name: "ShouldReturnNotFoundCode",
			httpErr: &error.HTTPError{
				StatusCode: 404,
				Code:       "not_found",
				Message:    "Not found",
			},
			expected: "not_found",
		},
		{
			name: "ShouldReturnValidationCode",
			httpErr: &error.HTTPError{
				StatusCode: 400,
				Code:       "validation_error",
				Message:    "Validation failed",
			},
			expected: "validation_error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := tt.httpErr.GetCode()
			suite.Equal(tt.expected, result)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestNewHTTPError() {
	tests := []struct {
		name       string
		statusCode int
		code       string
		message    string
		expected   *error.HTTPError
	}{
		{
			name:       "ShouldCreateNotFoundError",
			statusCode: 404,
			code:       "not_found",
			message:    "Resource not found",
			expected: &error.HTTPError{
				StatusCode: 404,
				Code:       "not_found",
				Message:    "Resource not found",
			},
		},
		{
			name:       "ShouldCreateValidationError",
			statusCode: 400,
			code:       "validation_error",
			message:    "Invalid input",
			expected: &error.HTTPError{
				StatusCode: 400,
				Code:       "validation_error",
				Message:    "Invalid input",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := error.NewHTTPError(tt.statusCode, tt.code, tt.message)
			suite.Equal(tt.expected.StatusCode, result.StatusCode)
			suite.Equal(tt.expected.Code, result.Code)
			suite.Equal(tt.expected.Message, result.Message)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestNewNotFoundError() {
	tests := []struct {
		name     string
		message  string
		expected *error.HTTPError
	}{
		{
			name:    "ShouldCreateNotFoundError",
			message: "User not found",
			expected: &error.HTTPError{
				StatusCode: 404,
				Code:       "not_found",
				Message:    "User not found",
			},
		},
		{
			name:    "ShouldCreateNotFoundErrorWithEmptyMessage",
			message: "",
			expected: &error.HTTPError{
				StatusCode: 404,
				Code:       "not_found",
				Message:    "",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := error.NewNotFoundError(tt.message)
			suite.Equal(tt.expected.StatusCode, result.StatusCode)
			suite.Equal(tt.expected.Code, result.Code)
			suite.Equal(tt.expected.Message, result.Message)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestNewValidationError() {
	tests := []struct {
		name     string
		message  string
		expected *error.HTTPError
	}{
		{
			name:    "ShouldCreateValidationError",
			message: "Email is required",
			expected: &error.HTTPError{
				StatusCode: 400,
				Code:       "validation_error",
				Message:    "Email is required",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := error.NewValidationError(tt.message)
			suite.Equal(tt.expected.StatusCode, result.StatusCode)
			suite.Equal(tt.expected.Code, result.Code)
			suite.Equal(tt.expected.Message, result.Message)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestNewUnauthorizedError() {
	tests := []struct {
		name     string
		message  string
		expected *error.HTTPError
	}{
		{
			name:    "ShouldCreateUnauthorizedError",
			message: "Invalid credentials",
			expected: &error.HTTPError{
				StatusCode: 401,
				Code:       "unauthorized",
				Message:    "Invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := error.NewUnauthorizedError(tt.message)
			suite.Equal(tt.expected.StatusCode, result.StatusCode)
			suite.Equal(tt.expected.Code, result.Code)
			suite.Equal(tt.expected.Message, result.Message)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestNewForbiddenError() {
	tests := []struct {
		name     string
		message  string
		expected *error.HTTPError
	}{
		{
			name:    "ShouldCreateForbiddenError",
			message: "Access denied",
			expected: &error.HTTPError{
				StatusCode: 403,
				Code:       "forbidden",
				Message:    "Access denied",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := error.NewForbiddenError(tt.message)
			suite.Equal(tt.expected.StatusCode, result.StatusCode)
			suite.Equal(tt.expected.Code, result.Code)
			suite.Equal(tt.expected.Message, result.Message)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestNewInternalServerError() {
	tests := []struct {
		name     string
		message  string
		expected *error.HTTPError
	}{
		{
			name:    "ShouldCreateInternalServerError",
			message: "Something went wrong",
			expected: &error.HTTPError{
				StatusCode: 500,
				Code:       "internal_server_error",
				Message:    "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := error.NewInternalServerError(tt.message)
			suite.Equal(tt.expected.StatusCode, result.StatusCode)
			suite.Equal(tt.expected.Code, result.Code)
			suite.Equal(tt.expected.Message, result.Message)
		})
	}
}

func (suite *HTTPErrorTestSuite) TestCommonErrorInstances() {
	suite.Run("ShouldHaveCorrectErrNotFound", func() {
		suite.Equal(404, error.ErrNotFound.GetStatusCode())
		suite.Equal("not_found", error.ErrNotFound.GetCode())
		suite.Equal("Data not found", error.ErrNotFound.Error())
	})

	suite.Run("ShouldHaveCorrectErrUnauthorized", func() {
		suite.Equal(401, error.ErrUnauthorized.GetStatusCode())
		suite.Equal("unauthorized", error.ErrUnauthorized.GetCode())
		suite.Equal("Unauthorized access", error.ErrUnauthorized.Error())
	})

	suite.Run("ShouldHaveCorrectErrForbidden", func() {
		suite.Equal(403, error.ErrForbidden.GetStatusCode())
		suite.Equal("forbidden", error.ErrForbidden.GetCode())
		suite.Equal("Access forbidden", error.ErrForbidden.Error())
	})
}
