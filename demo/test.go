package demo

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
func FromToTo(a From) To {
	var result To

	result.D = a.D
	result.E = a.E
	result.F = a.F
	result.A = a.A
	result.B = a.B

	return result
}
