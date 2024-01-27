package convert_time_to_and_from_numbers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapAB(t *testing.T) {
	orig := A{
		TimeToInt:   time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
		TimeToInt64: time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
	}
	dest := MapAB(orig)
	assert.Equal(t, int(orig.TimeToInt.Unix()), dest.TimeToInt)
	assert.Equal(t, orig.TimeToInt.Unix(), dest.TimeToInt64)
}
