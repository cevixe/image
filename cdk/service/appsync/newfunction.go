package appsync

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"
)

type FunctionProps struct {
	ApiId            string `field:"required"`
	DataSourceName   string `field:"required"`
	RequestTemplate  string `field:"required"`
	ResponseTemplate string `field:"required"`
}

func NewFunction(mod module.Module, alias string, props *FunctionProps) awsappsync.CfnFunctionConfiguration {

	name := naming.NewName(mod, naming.ResType_GraphQLFunction, alias)

	return awsappsync.NewCfnFunctionConfiguration(mod.Resource(), name.Logical(),
		&awsappsync.CfnFunctionConfigurationProps{
			ApiId:                   jsii.String(props.ApiId),
			DataSourceName:          jsii.String(props.DataSourceName),
			FunctionVersion:         jsii.String("2018-05-29"),
			Name:                    jsii.String(alias),
			RequestMappingTemplate:  jsii.String(props.RequestTemplate),
			ResponseMappingTemplate: jsii.String(props.ResponseTemplate),
		},
	)
}
