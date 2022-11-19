package lambda

import (
	"fmt"

	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	awsgo "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

func NewGolangFunction(mod module.Module, alias string, entry string) awslambda.Function {

	name := naming.NewName(mod, naming.ResType_Lambda, alias)
	role := NewFunctionRole(mod, alias)

	entryPath := fmt.Sprintf("%s/%s", mod.Location(), entry)

	return awsgo.NewGoFunction(mod.Resource(), name.Logical(), &awsgo.GoFunctionProps{
		FunctionName: name.Physical(),
		Architecture: awslambda.Architecture_X86_64(),
		Tracing:      awslambda.Tracing_ACTIVE,
		MemorySize:   jsii.Number(256),
		Entry:        jsii.String(entryPath),
		Runtime:      awslambda.Runtime_GO_1_X(),
		LogRetention: awslogs.RetentionDays_ONE_MONTH,
		Role:         role,
		Bundling: &awsgo.BundlingOptions{
			CgoEnabled: jsii.Bool(false),
			Environment: &map[string]*string{
				"GOOS":   jsii.String("linux"),
				"GOARCH": jsii.String("amd64"),
			},
			GoBuildFlags: &[]*string{
				jsii.String("-buildvcs=false"),
				jsii.String("-ldflags=\"-s -w\""),
			},
		},
	})
}
