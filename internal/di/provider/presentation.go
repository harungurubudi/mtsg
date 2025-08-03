package provider

import (
	"github.com/harungurubudi/mtsg/internal/presentation/http"
	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/internal/presentation/http/middleware"
	"github.com/harungurubudi/mtsg/internal/usecase"
	"github.com/harungurubudi/mtsg/pkg/config"
)

// Handler Providers

// ProvideHandlers provides HTTP handlers
func ProvideHandlers(authUseCase usecase.Authentication) *handler.Handlers {
	return handler.NewHandlers(authUseCase)
}

// Middleware Providers

// ProvideMiddlewareFactory provides middleware factory instance
func ProvideMiddlewareFactory(authUseCase usecase.Authentication) *middleware.Factory {
	return middleware.NewFactory(authUseCase)
}

// Server Providers

// ProvideHTTPServer provides HTTP server instance
func ProvideHTTPServer(
	handlers *handler.Handlers,
	config *config.Config,
	middlewareFactory *middleware.Factory,
) *http.Server {
	return http.NewServer(handlers, config, middlewareFactory)
}
