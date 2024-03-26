package uses_local_providers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"mapgenie/pkg/ptr"
	"testing"
)

func TestMapAB(t *testing.T) {
	from := A{
		V:  []byte("hello"),
		E:  fmt.Errorf("oops"),
		S1: nil,
		S2: ptr.Take("S2"),
	}
	to := MapAB(from)
	assert.Equal(t, to.V, int('h'))
	assert.Equal(t, to.E, []byte("oops"))
	assert.Equal(
		t, to.S1, "default",
		"expected provider mapper to take priority over value mapepr",
	)
	assert.Equal(t, to.S2, "S2")
}
