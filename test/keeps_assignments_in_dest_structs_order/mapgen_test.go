package keeps_assignments_in_dest_structs_order

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"regexp"
	"testing"
)

var assignRegex = regexp.MustCompile(`result\.(\w+)\s+=\s+input\.\w+`)

func checkAssignmentOrder(t *testing.T, path string, expectedOrder []string) {
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	matches := assignRegex.FindAllSubmatch(content, 100)

	assert.Len(t, matches, 3)
	for i, match := range matches {
		assert.Equal(t, expectedOrder[i], string(match[1]))
	}
}

func TestMapAB_DoesAssignmentsInOrder(t *testing.T) {
	checkAssignmentOrder(t, "./mapgen_ab.go", []string{"B", "C", "A"})
}

func TestMapBA(t *testing.T) {
	checkAssignmentOrder(t, "./mapgen_ba.go", []string{"A", "B", "C"})
}
