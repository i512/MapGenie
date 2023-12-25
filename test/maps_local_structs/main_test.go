package maps_local_structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptr[V any](v V) *V {
	return &v
}

func newA() A {
	return A{
		Int:            1,
		IntPtr:         ptr(2),
		UncommonIntA:   3,
		Str:            "4",
		StrPtr:         ptr("5"),
		UncommonStrA:   "6",
		ByteSlice:      []byte("7"),
		IntSlice:       []int{8},
		UncommonSliceA: []int{9},
		MapIntInt:      map[int]int{10: 10},
		ChanInt:        make(chan int),
		UncommonMapA:   map[byte]byte{},
		UncommonChanA:  make(chan string),
	}
}

func newB() B {
	return B{
		Int:            1,
		IntPtr:         ptr(2),
		UncommonIntB:   3,
		Str:            "4",
		StrPtr:         ptr("5"),
		UncommonStrB:   "6",
		ByteSlice:      []byte("7"),
		IntSlice:       []int{8},
		UncommonSliceB: []int{9},
		MapIntInt:      map[int]int{10: 10},
		ChanInt:        make(chan int),
		UncommonMapB:   map[byte]byte{},
		UncommonChanB:  make(chan string),
	}
}

func TestMapAA(t *testing.T) {
	origin := newA()
	dest := MapAA(origin)

	assert.Equal(t, origin.Int, dest.Int)
	assert.Equal(t, origin.IntPtr, dest.IntPtr)
	assert.Equal(t, origin.UncommonIntA, dest.UncommonIntA)
	assert.Equal(t, origin.Str, dest.Str)
	assert.Equal(t, origin.StrPtr, dest.StrPtr)
	assert.Equal(t, origin.UncommonStrA, origin.UncommonStrA)
	assert.Equal(t, origin.ByteSlice, dest.ByteSlice)
	assert.Equal(t, origin.IntSlice, dest.IntSlice)
	assert.Equal(t, origin.UncommonSliceA, dest.UncommonSliceA)
	assert.Equal(t, origin.MapIntInt, dest.MapIntInt)
	assert.Equal(t, origin.ChanInt, dest.ChanInt)
	assert.Equal(t, origin.UncommonMapA, origin.UncommonMapA)
	assert.Equal(t, origin.UncommonChanA, dest.UncommonChanA)
}

func TestMapAB(t *testing.T) {
	origin := newA()
	dest := MapAB(origin)

	assert.Equal(t, origin.Int, dest.Int)
	assert.Equal(t, origin.IntPtr, dest.IntPtr)
	assert.Equal(t, origin.Str, dest.Str)
	assert.Equal(t, origin.StrPtr, dest.StrPtr)
	assert.Equal(t, origin.ByteSlice, dest.ByteSlice)
	assert.Equal(t, origin.IntSlice, dest.IntSlice)
	assert.Equal(t, origin.MapIntInt, dest.MapIntInt)
	assert.Equal(t, origin.ChanInt, dest.ChanInt)
}

func TestMapBA(t *testing.T) {
	origin := newB()
	dest := MapBA(origin)

	assert.Equal(t, origin.Int, dest.Int)
	assert.Equal(t, origin.IntPtr, dest.IntPtr)
	assert.Equal(t, origin.Str, dest.Str)
	assert.Equal(t, origin.StrPtr, dest.StrPtr)
	assert.Equal(t, origin.ByteSlice, dest.ByteSlice)
	assert.Equal(t, origin.IntSlice, dest.IntSlice)
	assert.Equal(t, origin.MapIntInt, dest.MapIntInt)
	assert.Equal(t, origin.ChanInt, dest.ChanInt)
}
