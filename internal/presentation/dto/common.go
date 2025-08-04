package dto

import (
	"time"

	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
)

// ErrorResponse represents the standard error response format
// @Description Standard error response format
type ErrorResponse struct {
	Status string `json:"status" example:"error" description:"Response status"`
	Error  Error  `json:"error" description:"Error details"`
}

// Error contains error details
// @Description Error information
type Error struct {
	Code       string `json:"code" example:"validation_error" description:"Machine-readable error code"`
	Message    string `json:"message" example:"Invalid request format" description:"Human-readable error message"`
	StatusCode int    `json:"status_code" example:"400" description:"HTTP status code"`
}

// NewErrorResponse creates a new error response from HTTPError
func NewErrorResponse(httpError *pkgerror.HTTPError) *ErrorResponse {
	return &ErrorResponse{
		Status: "error",
		Error: Error{
			Code:       httpError.GetCode(),
			Message:    httpError.Error(),
			StatusCode: httpError.GetStatusCode(),
		},
	}
}

// SuccessResponse represents the standard success response format
// @Description Standard success response format
type SuccessResponse struct {
	Status string      `json:"status" example:"success" description:"Response status"`
	Data   interface{} `json:"data" description:"Response data"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status: "success",
		Data:   data,
	}
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Meta   Meta        `json:"meta"`
}

// Meta contains pagination metadata
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewListResponse creates a new list response with pagination
func NewListResponse(data interface{}, page, perPage, total int) *ListResponse {
	totalPages := (total + perPage - 1) / perPage
	return &ListResponse{
		Status: "success",
		Data:   data,
		Meta: Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page    int `json:"page" query:"page" validate:"min=1"`
	PerPage int `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

// GetOffset returns the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns the limit for database queries
func (p *PaginationRequest) GetLimit() int {
	return p.PerPage
}

// DefaultPaginationRequest returns default pagination values
func DefaultPaginationRequest() *PaginationRequest {
	return &PaginationRequest{
		Page:    1,
		PerPage: 20,
	}
}

// Timestamp represents a time value in ISO format
type Timestamp struct {
	time.Time
}

// MarshalJSON implements json.Marshaler
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if len(str) >= 2 {
		str = str[1 : len(str)-1]
	}

	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = parsed
	return nil
}
