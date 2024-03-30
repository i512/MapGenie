package resolves_name_collisions_when_importing_packages

import (
	"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/aliased"
	funccollision2 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/funccollision"
	"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/imported"
	typecollision2 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/typecollision"
	varcollision2 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/varcollision"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestMapAB(t *testing.T) {
	orig := A{
		T1: 1,
		T2: 2,
		T3: 3,
		T4: 4,
		T5: 5,
	}
	res := MapAB(orig)
	assert.Equal(t, funccollision2.T(1), res.T1)
	assert.Equal(t, varcollision2.T(2), res.T2)
	assert.Equal(t, typecollision2.T(3), res.T3)
	assert.Equal(t, imported.T(4), res.T4)
	assert.Equal(t, aliased.T(5), res.T5)
}

func TestUsesExistingImports(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "mapgen.go", nil, 0)
	require.NoError(t, err)

	imports := make(map[string]string)
	for _, imp := range file.Imports {
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		}
		imports[strings.Trim(imp.Path.Value, `"`)] = name
	}

	assert.Equal(t, map[string]string{
		"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/aliased":  "aliased3",
		"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/imported": "",

		"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/funccollision": "funccollision2",
		"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/varcollision":  "varcollision2",
		"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/typecollision": "typecollision2",
	}, imports)
}
