package middleware

import (
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// SENTRY will handle the SENTRY middleware
func (m *GoMiddleware) SENTRY(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span := sentry.StartSpan(
			c.Request().Context(), "", sentry.TransactionName(c.Request().URL.String()),
		)
		span.Finish()
		return next(c)
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
