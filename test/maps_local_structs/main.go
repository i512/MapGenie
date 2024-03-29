package maps_local_structs

type A struct {
	Int          int
	IntPtr       *int
	UncommonIntA int

	Str          string
	StrPtr       *string
	UncommonStrA string

	ByteSlice      []byte
	IntSlice       []int
	UncommonSliceA []int

	MapIntInt     map[int]int
	ChanInt       chan int
	UncommonMapA  map[byte]byte
	UncommonChanA chan string
}

type B struct {
	Int          int
	IntPtr       *int
	UncommonIntB int

	Str          string
	StrPtr       *string
	UncommonStrB string

	ByteSlice      []byte
	IntSlice       []int
	UncommonSliceB []int

	MapIntInt     map[int]int
	ChanInt       chan int
	UncommonMapB  map[byte]byte
	UncommonChanB chan string
}
