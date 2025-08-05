package http

import (
	"context"
	"net/http"
	"time"

	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/internal/presentation/http/middleware"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Server represents the HTTP server
type Server struct {
	echo     *echo.Echo
	handlers *handler.Handlers
	config   *config.Config
}

// NewServer creates a new HTTP server instance
func NewServer(handlers *handler.Handlers, config *config.Config) *Server {
	e := echo.New()

	// Configure Echo
	e.HideBanner = true
	e.HidePort = true

	// Set timeouts
	e.Server.ReadTimeout = config.Server.ReadTimeout
	e.Server.WriteTimeout = config.Server.WriteTimeout
	e.Server.IdleTimeout = config.Server.IdleTimeout

	// Apply global middleware
	e.Use(middleware.CreateGlobalMiddleware()...)

	return &Server{
		echo:     e,
		handlers: handlers,
		config:   config}
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Swagger documentation routes
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Ping endpoints (root level)
	s.echo.GET("/ping", s.handlers.Ping.Ping)
	s.echo.GET("/ping/protected", s.handlers.Ping.ProtectedPing, middleware.CreateAuthMiddleware())

	// Health check endpoint (root level)
	s.echo.GET("/health", s.healthCheck)

	// Debug endpoint to test container
	s.echo.GET("/debug/container", s.debugContainer)

	// API v1 routes
	v1 := s.echo.Group("/api/v1")

	// Public routes (no auth required)
	s.handlers.SetupRoutes(v1)
}

// healthCheck handles the health check endpoint
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// debugContainer handles the debug endpoint to test container injection
func (s *Server) debugContainer(c echo.Context) error {
	container, exists := middleware.GetContainerFromContext(c)
	if !exists {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Container not found in context",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"container_found": true,
		"authentication":  container.Authentication != nil,
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()

	s.echo.Logger.Infof("Starting server on port %s", s.config.Server.Port)
	return s.echo.Start(":" + s.config.Server.Port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.echo.Logger.Info("Shutting down server...")
	return s.echo.Shutdown(ctx)
}

// NewConfig creates a new server configuration from environment variables
// This function is deprecated and should be replaced with pkg/config.Load()
func NewConfig() *config.Config {
	// This is a temporary function for backward compatibility
	// In the future, use pkg/config.Load() directly
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
