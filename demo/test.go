package demo

import (
	"mopper/demo/demo2"
	demo3aliased "mopper/demo/demo2/demo3"
)

type From struct {
	A	int
	B	string
	c	int
	D	[]byte
	E	[]int
	F	float64
}

type To struct {
	Other	int
	A	int
	B	string
	D	[]byte
	E	[]int
	F	float64
}

// Map1 map this pls
func Map1(a From) To {
	var result To

	result.A = a.A
	result.B = a.B
	result.D = a.D
	result.E = a.E
	result.F = a.F

	return result
}

// Map2 map this pls
func Map2(a From) demo2.Demo2S {
	var result demo2.Demo2S

	result.D = a.D
	result.E = a.E
	result.F = a.F
	result.A = a.A
	result.B = a.B

	return result
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
