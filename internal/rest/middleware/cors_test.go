package middleware_test

import (
	"net/http"
	test "net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/josemolinai-nuamx/go-clean-arch/internal/rest/middleware"
)

func TestCORS(t *testing.T) {
	t.Run("debug-fallback-allows-any-origin", func(t *testing.T) {
		t.Setenv("DEBUG", "true")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")

		e := echo.New()
		req := test.NewRequest(echo.GET, "/", nil)
		res := test.NewRecorder()
		c := e.NewContext(req, res)

		h := middleware.CORS(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))

		err := h(c)
		require.NoError(t, err)
		assert.Equal(t, "*", res.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("allows-configured-origin", func(t *testing.T) {
		t.Setenv("DEBUG", "false")
		t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000, http://localhost:5173")

		e := echo.New()
		req := test.NewRequest(echo.GET, "/", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		res := test.NewRecorder()
		c := e.NewContext(req, res)

		h := middleware.CORS(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))

		err := h(c)
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:3000", res.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("blocks-origin-not-in-allowlist", func(t *testing.T) {
		t.Setenv("DEBUG", "false")
		t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

		e := echo.New()
		req := test.NewRequest(echo.GET, "/", nil)
		req.Header.Set("Origin", "https://evil.example")
		res := test.NewRecorder()
		c := e.NewContext(req, res)

		h := middleware.CORS(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))

		err := h(c)
		require.NoError(t, err)
		assert.Empty(t, res.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestMain(m *testing.M) {
	os.Setenv("DEBUG", "")
	os.Setenv("CORS_ALLOWED_ORIGINS", "")
	os.Exit(m.Run())
}
