package middleware_test

import (
	"net/http"
	test "net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/josemolinai-nuamx/go-clean-arch/internal/rest/middleware"
)

func TestSetRequestContextWithTimeout(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewRecorder()
	c := e.NewContext(req, res)

	var requestDeadline time.Time
	var hasDeadline bool
	var doneBeforeReturn bool
	var middlewareContextDone <-chan struct{}

	h := middleware.SetRequestContextWithTimeout(50 * time.Millisecond)(echo.HandlerFunc(func(c echo.Context) error {
		ctx := c.Request().Context()
		requestDeadline, hasDeadline = ctx.Deadline()
		select {
		case <-ctx.Done():
			doneBeforeReturn = true
		default:
			doneBeforeReturn = false
		}
		middlewareContextDone = ctx.Done()
		return c.NoContent(http.StatusOK)
	}))

	err := h(c)
	require.NoError(t, err)

	assert.True(t, hasDeadline)
	assert.False(t, doneBeforeReturn)
	assert.Equal(t, http.StatusOK, res.Code)

	remaining := time.Until(requestDeadline)
	assert.LessOrEqual(t, remaining, 50*time.Millisecond)
	assert.Greater(t, remaining, 0*time.Millisecond)

	select {
	case <-middlewareContextDone:
		assert.True(t, true)
	default:
		t.Fatal("expected context to be canceled after middleware returns")
	}
}
