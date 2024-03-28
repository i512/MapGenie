package uses_local_providers

type A struct {
	V  []byte
	E  error
	S1 *string
	S2 *string
}

type B struct {
	V  int
	E  []byte
	S1 string
	S2 string
}

// BytesToInt magic provider
func BytesToInt(b []byte) int {
	return int(b[0])
}

// StrValueWithDefault magic provider
func StrValueWithDefault(s *string) string {
	if s == nil {
		return "default"
	}

	return *s
}
