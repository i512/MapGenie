package demo

import (
	"mopper/demo/demo2"
	demo3aliased "mopper/demo/demo2/demo3"
)

type From struct {
	A int
	B string
	c int
	D []byte
	E []int
	F float64
}

type To struct {
	Other int
	A     int
	B     string
	D     []byte
	E     []int
	F     float64
}

// Map1 map this pls
func Map1(a From) To {
	return To{}
}

// Map2 map this pls
func Map2(a From) demo2.Demo2S {
	return demo2.Demo2S{}
}

// Map3 map this pls
func Map3(a From) demo3aliased.Demo3S {
	var result demo3aliased.Demo3S

	result.A = a.A
	result.B = a.B
	result.D = a.D
	result.E = a.E
	result.F = a.F

	return result
}
