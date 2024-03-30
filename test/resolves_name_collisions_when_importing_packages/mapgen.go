package resolves_name_collisions_when_importing_packages

import (
	aliased3 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/aliased"
	funccollision2 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/funccollision"
	"github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/imported"
	typecollision2 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/typecollision"
	varcollision2 "github.com/i512/mapgenie/test/resolves_name_collisions_when_importing_packages/varcollision"
)

// MapAB map this pls
func MapAB(input A) B {
	return B{T1: funccollision2.T(input.T1), T2: varcollision2.T(input.T2), T3: typecollision2.T(input.T3), T4: imported.T(input.T4), T5: aliased3.T(input.T5)}
}

// MapBA map this pls
func MapBA(input B) A {
	return A{T2: int(input.T2), T3: int(input.T3), T4: int(input.T4), T5: int(input.T5), T1: int(input.T1)}
}
