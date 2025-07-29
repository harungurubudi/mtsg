package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server represents the HTTP server
type Server struct {
	echo     *echo.Echo
	handlers *handler.Handlers
	config   *Config
}

// Config represents server configuration
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewServer creates a new HTTP server instance
func NewServer(handlers *handler.Handlers, config *Config) *Server {
	e := echo.New()

	// Configure Echo
	e.HideBanner = true
	e.HidePort = true

	// Set timeouts
	e.Server.ReadTimeout = config.ReadTimeout
	e.Server.WriteTimeout = config.WriteTimeout
	e.Server.IdleTimeout = config.IdleTimeout

	return &Server{
		echo:     e,
		handlers: handlers,
		config:   config,
	}
}

// setupMiddleware configures global middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware for panic handling
	s.echo.Use(middleware.Recover())

	// Logger middleware for request logging
	s.echo.Use(middleware.Logger())

	// CORS middleware for cross-origin requests
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Custom middleware
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.Gzip())
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Ping endpoint
	s.echo.GET("/ping", s.handlers.Ping.Ping)

	// Health check endpoint
	s.echo.GET("/health", s.healthCheck)

	// API v1 routes (for future endpoints)
	v1 := s.echo.Group("/api/v1")

	// Public routes (no auth required)
	_ = v1.Group("")
	// TODO: Add authentication endpoints here
	// public.POST("/auth/login", s.handlers.Auth.Login)

	// Protected routes (auth required)
	_ = v1.Group("")
	// TODO: Add authentication middleware here
	// protected.Use(AuthMiddleware(s.handlers.Auth.UseCase))
	// protected.GET("/users/profile", s.handlers.User.GetProfile)
}

// healthCheck handles the health check endpoint
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// setupErrorHandler configures custom error handling
func (s *Server) setupErrorHandler() {
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			// Echo HTTP errors
			c.JSON(httpErr.Code, map[string]interface{}{
				"error":   httpErr.Message,
				"code":    httpErr.Code,
				"details": httpErr.Internal,
			})
			return
		}

		// Internal server errors
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
			"code":  "INTERNAL_ERROR",
		})
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupMiddleware()
	s.setupErrorHandler()
	s.setupRoutes()

	s.echo.Logger.Infof("Starting server on port %s", s.config.Port)
	return s.echo.Start(":" + s.config.Port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.echo.Logger.Info("Shutting down server...")
	return s.echo.Shutdown(ctx)
}

// NewConfig creates a new server configuration from environment variables
func NewConfig() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	readTimeout, _ := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if readTimeout == 0 {
		readTimeout = 30 * time.Second
	}

	writeTimeout, _ := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if writeTimeout == 0 {
		writeTimeout = 30 * time.Second
	}

	idleTimeout, _ := time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if idleTimeout == 0 {
		idleTimeout = 60 * time.Second
	}

	return &Config{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}
