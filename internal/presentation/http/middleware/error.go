package middleware

import (
	"net/http"
	"strconv"

	"github.com/harungurubudi/mtsg/internal/presentation/dto"
	pkgerror "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/labstack/echo/v4"
)

// HttpErrorHandler provides centralized error handling middleware
// It handles different error types and ensures consistent error response format
func HttpErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				// Handle custom HTTP errors (pkgerror.HTTPError)
				if httpError, ok := err.(*pkgerror.HTTPError); ok {
					// Log the actual error for debugging while returning sanitized response
					c.Logger().Errorf("HTTP Error: %v", httpError)

					if httpError.GetStatusCode() == http.StatusInternalServerError {
						return c.JSON(httpError.GetStatusCode(), dto.NewErrorResponse(pkgerror.NewInternalServerError("internal server error")))
					}
					return c.JSON(httpError.GetStatusCode(), dto.NewErrorResponse(httpError))
				}

				// Handle Echo framework HTTP errors
				if httpError, ok := err.(*echo.HTTPError); ok {
					c.Logger().Errorf("Echo HTTP Error: %v", httpError)

					return c.JSON(httpError.Code, &dto.ErrorResponse{
						Status: "error",
						Error: dto.Error{
							Code:       strconv.Itoa(httpError.Code),
							Message:    httpError.Error(),
							StatusCode: httpError.Code,
						},
					})
				}

				// Handle generic/unexpected errors
				c.Logger().Errorf("Unexpected error: %v", err)
				return c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(pkgerror.NewInternalServerError("internal server error")))
			}

			return nil
		}
	}
}
