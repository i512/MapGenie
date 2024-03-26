package uses_local_providers

import "fmt"

// ErrorToByteSlice magic provider
func ErrorToByteSlice(e error) ([]byte, error) {
	if e == nil {
		return nil, fmt.Errorf("oh no!")
	}

	return []byte(e.Error()), nil
}
