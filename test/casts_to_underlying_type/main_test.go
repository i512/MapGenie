package casts_to_underlying_type

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{A: "path"}
	dest := MapAB(orig)

	assert.Equal(t, URL("path"), dest)
}

func TestMapBA(t *testing.T) {
	orig := B{A: URL("path")}
	dest := MapBA(orig)

	assert.Equal(t, "path", dest)
}
