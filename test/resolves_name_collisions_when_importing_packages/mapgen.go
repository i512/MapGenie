package resolves_name_collisions_when_importing_packages

import (
	aliased3 "mapgenie/test/resolves_name_collisions_when_importing_packages/aliased"
	"mapgenie/test/resolves_name_collisions_when_importing_packages/imported"
)

// MapAB map this pls
func MapAB(A) B {
	return B{
		T4: imported.T(-1),
		T5: aliased3.T(-1),
	}
}

// MapBA map this pls
func MapBA(B) A {
	return A{}
}
