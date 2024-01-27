package convert_time_to_and_from_string

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapAB(t *testing.T) {
	tt := time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)
	orig := A{T: tt}
	dest := MapAB(orig)
	assert.Equal(t, tt.Format(time.RFC3339), dest.T)
}
