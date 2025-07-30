package errorhandler

import (
	"github.com/harungurubudi/mtsg/pkg/config"
)

// ErrorHandlerFactory provides factory functions for error handlers
type ErrorHandlerFactory struct {
	config *config.Config
}

// NewErrorHandlerFactory creates a new error handler factory
func NewErrorHandlerFactory(config *config.Config) *ErrorHandlerFactory {
	return &ErrorHandlerFactory{
		config: config,
	}
}

// CreateErrorHandler creates an error handler with configuration
func (f *ErrorHandlerFactory) CreateErrorHandler() *ErrorHandler {
	errorConfig := &Config{
		Environment: f.config.Server.Environment,
		ShowDetails: f.config.Server.Environment == "development",
		LogErrors:   true,
	}

	return NewErrorHandler(errorConfig)
}

// CreateErrorHandlerWithConfig creates an error handler with custom configuration
func (f *ErrorHandlerFactory) CreateErrorHandlerWithConfig(errorConfig *Config) *ErrorHandler {
	return NewErrorHandler(errorConfig)
}

// CreateDevelopmentErrorHandler creates an error handler for development environment
func (f *ErrorHandlerFactory) CreateDevelopmentErrorHandler() *ErrorHandler {
	errorConfig := &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   true,
	}

	return NewErrorHandler(errorConfig)
}

// CreateProductionErrorHandler creates an error handler for production environment
func (f *ErrorHandlerFactory) CreateProductionErrorHandler() *ErrorHandler {
	errorConfig := &Config{
		Environment: "production",
		ShowDetails: false,
		LogErrors:   true,
	}

	return NewErrorHandler(errorConfig)
}
