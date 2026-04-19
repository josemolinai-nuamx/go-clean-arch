package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/josemolinai-nuamx/go-clean-arch/domain"
)

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{name: "nil", err: nil, expected: http.StatusOK},
		{name: "internal-server-error", err: domain.ErrInternalServerError, expected: http.StatusInternalServerError},
		{name: "not-found", err: domain.ErrNotFound, expected: http.StatusNotFound},
		{name: "conflict", err: domain.ErrConflict, expected: http.StatusConflict},
		{name: "unknown", err: errors.New("unknown"), expected: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, getStatusCode(tt.err))
		})
	}
}
