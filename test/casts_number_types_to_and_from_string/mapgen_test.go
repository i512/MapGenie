package casts_number_types_to_and_from_string

import (
	"fmt"
	"github.com/i512/mapgenie/test/casts_number_types_to_and_from_string/custom_fmt_included"
	"github.com/i512/mapgenie/test/casts_number_types_to_and_from_string/fmt_aliased"
	"github.com/i512/mapgenie/test/casts_number_types_to_and_from_string/fmt_defined"
	"github.com/i512/mapgenie/test/casts_number_types_to_and_from_string/fmt_included"
	"github.com/i512/mapgenie/test/casts_number_types_to_and_from_string/structs"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func checkFormatsNumber(t *testing.T, f func(structs.A) structs.B) {
	orig := structs.A{
		Int:     1,
		Int8:    2,
		Int16:   3,
		Int32:   4,
		Int64:   5,
		Uint:    6,
		Uint8:   7,
		Uint16:  8,
		Uint32:  9,
		Uint64:  10,
		Float64: 11.1,
		Float32: 12.1,
	}
	dest := f(orig)

	assert.Equal(t, "1", dest.Int)
	assert.Equal(t, "2", dest.Int8)
	assert.Equal(t, "3", dest.Int16)
	assert.Equal(t, "4", dest.Int32)
	assert.Equal(t, "5", dest.Int64)
	assert.Equal(t, "6", dest.Uint)
	assert.Equal(t, "7", dest.Uint8)
	assert.Equal(t, "8", dest.Uint16)
	assert.Equal(t, "9", dest.Uint32)
	assert.Equal(t, "10", dest.Uint64)
	assert.Equal(t, "11.1", dest.Float64)
	assert.Equal(t, "12.1", dest.Float32)

}

func checkParsesNumber(t *testing.T, f func(structs.B) structs.A) {
	orig := structs.B{
		Int:     strconv.FormatInt(-1<<63, 10),
		Int8:    strconv.FormatInt(-1<<7, 10),
		Int16:   strconv.FormatInt(-1<<15, 10),
		Int32:   strconv.FormatInt(-1<<31, 10),
		Int64:   strconv.FormatInt(-1<<63, 10),
		Uint:    strconv.FormatUint(1<<64-1, 10),
		Uint8:   strconv.FormatUint(1<<8-1, 10),
		Uint16:  strconv.FormatUint(1<<16-1, 10),
		Uint32:  strconv.FormatUint(1<<32-1, 10),
		Uint64:  strconv.FormatUint(1<<64-1, 10),
		Float64: fmt.Sprint("6.4"),
		Float32: fmt.Sprint("3.2"),
	}

	dest := f(orig)
	assert.Equal(t, int(-1<<63), dest.Int)
	assert.Equal(t, int8(-1<<7), dest.Int8)
	assert.Equal(t, int16(-1<<15), dest.Int16)
	assert.Equal(t, int32(-1<<31), dest.Int32)
	assert.Equal(t, int64(-1<<63), dest.Int64)
	assert.Equal(t, uint(1<<64-1), dest.Uint)
	assert.Equal(t, uint8(1<<8-1), dest.Uint8)
	assert.Equal(t, uint16(1<<16-1), dest.Uint16)
	assert.Equal(t, uint32(1<<32-1), dest.Uint32)
	assert.Equal(t, uint64(1<<64-1), dest.Uint64)
	assert.Equal(t, float64(6.4), dest.Float64)
	assert.Equal(t, float32(3.2), dest.Float32)
}

func TestMapAB(t *testing.T) {
	checkFormatsNumber(t, MapAB)
	checkParsesNumber(t, MapBA)
}

func TestMapABCustomFmtIncluded(t *testing.T) {
	checkFormatsNumber(t, custom_fmt_included.MapAB)
	checkParsesNumber(t, custom_fmt_included.MapBA)
}

func TestMapABFmtAliased(t *testing.T) {
	checkFormatsNumber(t, fmt_aliased.MapAB)
	checkParsesNumber(t, fmt_aliased.MapBA)
}

func TestMapABFmtDefined(t *testing.T) {
	checkFormatsNumber(t, fmt_defined.MapAB)
	checkParsesNumber(t, fmt_defined.MapBA)
}

func TestMapABFmtIncluded(t *testing.T) {
	checkFormatsNumber(t, fmt_included.MapAB)
	checkParsesNumber(t, fmt_included.MapBA)
}
