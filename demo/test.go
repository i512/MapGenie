package demo

import "mopper/demo/demo2"

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

// FromToTo map this pls
func FromToTo(a From) demo2.S2 {
	var result demo2.S2

	result.A = a.A
	result.B = a.B
	result.D = a.D
	result.E = a.E
	result.F = a.F

	return result
}
