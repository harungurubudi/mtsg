package http

import (
	"context"
	"net/http"
	"time"

	"github.com/harungurubudi/mtsg/internal/presentation/http/errorhandler"
	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/pkg/config"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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

	return &Server{
		echo:     e,
		handlers: handlers,
		config:   config,
	}
}

// setupMiddleware configures global middleware
func (s *Server) setupMiddleware() {
	// TODO: Create middleware factory with proper dependencies
	// For now, use basic Echo middleware
	s.echo.Use(echoMiddleware.Recover())
	s.echo.Use(echoMiddleware.Logger())
	s.echo.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	s.echo.Use(echoMiddleware.RequestID())
	s.echo.Use(echoMiddleware.Gzip())
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Swagger documentation routes
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Ping endpoints
	s.echo.GET("/ping", s.handlers.Ping.Ping)

	// Health check endpoint
	s.echo.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.echo.Group("/api/v1")

	// Public routes (no auth required)
	public := v1.Group("")
	public.POST("/auth/login", s.handlers.Auth.Login)
	public.POST("/auth/refresh", s.handlers.Auth.RefreshToken)

	// Protected routes (auth required)
	protected := v1.Group("")
	// TODO: Add authentication middleware here
	// protected.Use(AuthMiddleware(s.handlers.Auth.UseCase))
	protected.POST("/auth/logout", s.handlers.Auth.Logout)
	protected.GET("/ping/protected", s.handlers.Ping.ProtectedPing)
	// TODO: Add other protected endpoints here
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
	factory := errorhandler.NewErrorHandlerFactory(s.config)
	errorHandler := factory.CreateErrorHandler()

	s.echo.HTTPErrorHandler = errorHandler.HandleError
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupMiddleware()
	s.setupErrorHandler()
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
