package keeps_assignments_in_dest_structs_order

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"regexp"
	"testing"
)

func TestMapAB_DoesAssignmentsInOrder(t *testing.T) {
	content, err := os.ReadFile("./mapgen_ab.go")
	require.NoError(t, err)

	regex := regexp.MustCompile(`result\.(\w+)\s+=\s+input\.\w+`)
	matches := regex.FindAll(content, 100)

	expectedOrder := []string{"B", "C", "A"}

	assert.Len(t, matches, 3)
	for i, match := range matches {
		assert.Equal(t, expectedOrder[i], match[1])
	}
}

func TestMapBA(t *testing.T) {
}
