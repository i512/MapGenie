package resolves_name_collisions_when_importing_packages

import (
	"mapgenie/test/resolves_name_collisions_when_importing_packages/aliased"
	funccollision2 "mapgenie/test/resolves_name_collisions_when_importing_packages/funccollision"
	"mapgenie/test/resolves_name_collisions_when_importing_packages/imported"
	typecollision2 "mapgenie/test/resolves_name_collisions_when_importing_packages/typecollision"
	varcollision2 "mapgenie/test/resolves_name_collisions_when_importing_packages/varcollision"
)

type A struct {
	T1 int
	T2 int
	T3 int
	T4 int
	T5 int
}

type B struct {
	T1 funccollision2.T
	T2 varcollision2.T
	T3 typecollision2.T
	T4 imported.T
	T5 aliased.T
}

func funccollision() {}

var varcollision int

type typecollision struct{}
