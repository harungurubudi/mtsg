package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harungurubudi/mtsg/internal/presentation/http/handler"
	"github.com/harungurubudi/mtsg/pkg/config"
)

func TestSwaggerEndpoints(t *testing.T) {
	// Create test server
	mockAuth := new(MockAuthentication)
	handlers := handler.NewHandlers(mockAuth)
	config := &config.Config{Server: config.ServerConfig{Port: "8080"}}
	server := NewServer(handlers, config)

	// Setup routes
	server.setupRoutes()

	// Test Swagger UI endpoint
	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rec := httptest.NewRecorder()
	c := server.echo.NewContext(req, rec)

	// Find the route
	server.echo.Router().Find(http.MethodGet, "/swagger/index.html", c)

	// Test Swagger JSON endpoint
	req = httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	rec = httptest.NewRecorder()
	c = server.echo.NewContext(req, rec)

	server.echo.Router().Find(http.MethodGet, "/swagger/doc.json", c)
}
