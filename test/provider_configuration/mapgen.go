package provider_configuration

// mapgenie providers:
// . github.com/i512/test/provider_configuration/default_providers
// github.com/i512/test/provider_configuration/providers
// mapgenie default mapping:
// RecordID => ID
// PostID => ID
// prefix(\w+)ID => ID$1

// MapAB map this pls:
// Field1 => OtherField1
// Conv1(Field2) => OtherField2
// Conv2(Field3)
// providers.Conv3(Field4)
func MapAB(A) B {
	return B{}
}
