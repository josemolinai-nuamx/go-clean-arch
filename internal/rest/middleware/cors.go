package middleware

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const defaultDevOrigin = "*"

func parseAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		if strings.EqualFold(os.Getenv("DEBUG"), "true") {
			return []string{defaultDevOrigin}
		}
		return nil
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func originAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == defaultDevOrigin || allowed == origin {
			return true
		}
	}

	return false
}

// CORS will handle the CORS middleware
func CORS(next echo.HandlerFunc) echo.HandlerFunc {
	allowedOrigins := parseAllowedOrigins()

	return func(c echo.Context) error {
		origin := strings.TrimSpace(c.Request().Header.Get("Origin"))
		if originAllowed(origin, allowedOrigins) {
			if origin == "" {
				c.Response().Header().Set("Access-Control-Allow-Origin", defaultDevOrigin)
			} else {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		return next(c)
	}
}
