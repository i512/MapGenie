package uses_local_providers

type A struct {
	V []byte
}

type B struct {
	V int
}

// BytesToInt magic provider
func BytesToInt(b []byte) int {
	return int(b[0])
}
