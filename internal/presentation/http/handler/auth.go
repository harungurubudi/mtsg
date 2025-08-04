package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/harungurubudi/mtsg/internal/domain/authentication"
	"github.com/harungurubudi/mtsg/internal/presentation/dto"
	"github.com/harungurubudi/mtsg/internal/presentation/http/middleware"
	http_error "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/harungurubudi/mtsg/pkg/token"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	validator *validator.Validate
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		validator: validator.New(),
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.SuccessResponse{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} error.HTTPError "Invalid request format or validation error"
// @Failure 401 {object} error.HTTPError "Invalid credentials"
// @Failure 404 {object} error.HTTPError "User not found"
// @Failure 403 {object} error.HTTPError "Account is inactive"
// @Failure 500 {object} error.HTTPError "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	container, exists := middleware.GetContainerFromContext(c)
	if !exists {
		return http_error.NewInternalServerError("DI container not available")
	}
	ctx := c.Request().Context()

	// Parse and validate request
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return http_error.NewValidationError("invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return http_error.NewValidationError("validation failed")
	}

	// Convert to domain model
	credential := req.ToDomain()

	// Call use case
	session, err := container.Authentication.Login(ctx, &credential)
	if err != nil {
		// Handle specific domain errors
		switch {
		case errors.Is(err, authentication.ErrInvalidCredential):
			return http_error.NewUnauthorizedError("invalid credentials")
		case errors.Is(err, authentication.ErrUserNotFound):
			return http_error.NewNotFoundError("user not found")
		case errors.Is(err, authentication.ErrUserInactive):
			return http_error.NewForbiddenError("account is inactive")
		default:
			return http_error.NewInternalServerError("login failed")
		}
	}

	// Return success response
	response := dto.NewLoginResponse(session)
	return c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

// Logout godoc
// @Summary User logout
// @Description Logout user by invalidating the access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.SuccessResponse{data=object} "Logout successful"
// @Failure 401 {object} error.HTTPError "Missing or invalid authorization header"
// @Failure 500 {object} error.HTTPError "Internal server error"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	container, exists := middleware.GetContainerFromContext(c)
	if !exists {
		return http_error.NewInternalServerError("DI container not available")
	}
	ctx := c.Request().Context()

	// Extract token from request using the extracted function
	tokenStr, err := middleware.ExtractTokenFromRequest(c)
	if err != nil {
		return err
	}

	// Call use case
	err = container.Authentication.Logout(ctx, tokenStr)
	if err != nil {
		return http_error.NewInternalServerError("logout failed")
	}

	return c.JSON(http.StatusOK, dto.NewSuccessResponse(map[string]string{
		"message": "successfully logged out",
	}))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access and refresh tokens using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.SuccessResponse{data=dto.LoginResponse} "Token refresh successful"
// @Failure 400 {object} error.HTTPError "Invalid request format or validation error"
// @Failure 401 {object} error.HTTPError "Invalid or expired refresh token"
// @Failure 500 {object} error.HTTPError "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	container, exists := middleware.GetContainerFromContext(c)
	if !exists {
		return http_error.NewInternalServerError("DI container not available")
	}
	ctx := c.Request().Context()

	// Parse request
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return http_error.NewValidationError("invalid request format")
	}

	if err := h.validator.Struct(req); err != nil {
		return http_error.NewValidationError("validation failed")
	}

	// Call use case
	session, err := container.Authentication.RefreshToken(ctx, token.Token(req.RefreshToken))
	if err != nil {
		switch {
		case errors.Is(err, authentication.ErrInvalidAuthentication):
			return http_error.NewUnauthorizedError("invalid refresh token")
		case errors.Is(err, authentication.ErrUserNotFound):
			return http_error.NewUnauthorizedError("refresh token expired")
		default:
			return http_error.NewInternalServerError("token refresh failed")
		}
	}

	// Return new tokens
	response := dto.NewLoginResponse(session)
	return c.JSON(http.StatusOK, dto.NewSuccessResponse(response))
}

func (h *AuthHandler) SetupRoutes(prefix string, e *echo.Group) {
	auth := e.Group(prefix)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.RefreshToken)
}
