package casts_number_types_to_string

import (
	"github.com/stretchr/testify/assert"
	"mapgenie/test/casts_number_types_to_string/custom_fmt_included"
	"mapgenie/test/casts_number_types_to_string/fmt_aliased"
	"mapgenie/test/casts_number_types_to_string/fmt_defined"
	"mapgenie/test/casts_number_types_to_string/fmt_included"
	"mapgenie/test/casts_number_types_to_string/structs"
	"testing"
)

func checkFunc(t *testing.T, f func(structs.A) structs.B) {
	orig := structs.A{
		Int:    1,
		Int8:   2,
		Int16:  3,
		Int32:  4,
		Int64:  5,
		Uint:   6,
		Uint8:  7,
		Uint16: 8,
		Uint32: 9,
		Uint64: 10,
	}
	dest := f(orig)

	assert.Equal(t, orig.Int, dest.Int)
	assert.Equal(t, orig.Uint, dest.Uint)
}

func TestMapAB(t *testing.T) {
	checkFunc(t, MapAB)
}

func TestMapABCustomFmtIncluded(t *testing.T) {
	checkFunc(t, custom_fmt_included.MapAB)
}

func TestMapABFmtAliased(t *testing.T) {
	checkFunc(t, fmt_aliased.MapAB)
}

func TestMapABFmtDefined(t *testing.T) {
	checkFunc(t, fmt_defined.MapAB)
}

func TestMapABFmtIncluded(t *testing.T) {
	checkFunc(t, fmt_included.MapAB)
}
