// Package magic implements magic comment analysis
package magic

//mapgenie providers:
//. github.com/i512/test/provider_configuration/default_providers
//github.com/i512/test/provider_configuration/providers
//mapgenie default mapping:
//RecordID => ID
//PostID => ID
//prefix(\w+)ID => ID$1

type Mapping struct {
	To             []string
	CustomProvider string
}

type ProviderPackage struct {
	Pkg   string
	Alias string
}

type LocalConfig interface {
	Mappings(s string) []Mapping
	ProviderPackages() []ProviderPackage
}

type LocalConfigH struct {
}

func ParseComment(s string) LocalConfig {

}
