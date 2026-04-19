package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeCursorRoundTrip(t *testing.T) {
	original := time.Date(2026, 4, 18, 12, 30, 45, 987000000, time.UTC)
	encoded := EncodeCursor(original)
	decoded, err := DecodeCursor(encoded)

	require.NoError(t, err)
	assert.Equal(t, original.Format(timeFormat), decoded.Format(timeFormat))
}

func TestDecodeCursorInvalidBase64(t *testing.T) {
	_, err := DecodeCursor("%%%")
	require.Error(t, err)
}
