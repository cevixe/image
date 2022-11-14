package module

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type impl struct {
	app      string
	name     string
	location string
	resource constructs.Construct
}

func (m *impl) App() string {
	return m.app
}

func (m *impl) Name() string {
	return m.name
}

func (m *impl) Location() string {
	return m.location
}

func (m *impl) Export(name string, value string) string {

	varname := fmt.Sprintf("%s-%s", m.name, name)
	awscdk.NewCfnOutput(m.Resource(), jsii.String(name), &awscdk.CfnOutputProps{
		Value:      jsii.String(value),
		ExportName: jsii.String(varname),
	})

	return varname
}

func (m *impl) Import(from string, name string) string {

	varname := fmt.Sprintf("%s-%s", from, name)
	externalValue := awscdk.Fn_ImportValue(jsii.String(varname))
	return *externalValue
}

func (m *impl) DependsOn(mod ...Module) {
	for _, item := range mod {
		m.resource.Node().AddDependency(item.Resource())
	}
}

func (m *impl) Resource() constructs.Construct {
	return m.resource
}
