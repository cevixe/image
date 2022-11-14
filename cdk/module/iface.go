package module

import "github.com/aws/constructs-go/constructs/v10"

type Module interface {
	App() string
	Name() string
	Location() string
	DependsOn(mods ...Module)
	Import(from string, name string) string
	Export(name string, value string) string
	Resource() constructs.Construct
}
