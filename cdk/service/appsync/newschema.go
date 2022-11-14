package appsync

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"
)

type SchemaProps struct {
	ApiId      string `field:"required"`
	Definition string `field:"required"`
}

func NewSchema(mod module.Module, alias string, props *SchemaProps) awsappsync.CfnGraphQLSchema {

	name := naming.NewName(mod, naming.ResType_GraphQLSchema, alias)

	return awsappsync.NewCfnGraphQLSchema(
		mod.Resource(),
		name.Logical(),
		&awsappsync.CfnGraphQLSchemaProps{
			ApiId:      jsii.String(props.ApiId),
			Definition: jsii.String(props.Definition),
		},
	)
}
