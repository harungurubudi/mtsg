package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoggingMiddleware provides request logging
func LoggingMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} id=${id} remote_ip=${remote_ip} " +
			"host=${host} method=${method} uri=${uri} " +
			"user_agent=${user_agent} status=${status} " +
			"error=${error} latency=${latency} latency_human=${latency_human} " +
			"bytes_in=${bytes_in} bytes_out=${bytes_out}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	})
}
