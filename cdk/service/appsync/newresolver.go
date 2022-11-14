package appsync

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"
)

type ResolverProps struct {
	ApiId            string   `field:"required"`
	Type             string   `field:"required"`
	Field            string   `field:"required"`
	Functions        []string `field:"required"`
	RequestTemplate  string   `field:"required"`
	ResponseTemplate string   `field:"required"`
}

func NewResolver(mod module.Module, alias string, props *ResolverProps) awsappsync.CfnResolver {

	name := naming.NewName(mod, naming.ResType_GraphQLResolver, alias)

	functions := make([]*string, 0)
	for _, fn := range props.Functions {
		functions = append(functions, jsii.String(fn))
	}

	return awsappsync.NewCfnResolver(
		mod.Resource(),
		name.Logical(),
		&awsappsync.CfnResolverProps{
			ApiId:                   jsii.String(props.ApiId),
			TypeName:                jsii.String(props.Type),
			FieldName:               jsii.String(props.Field),
			Kind:                    jsii.String("PIPELINE"),
			RequestMappingTemplate:  jsii.String(props.RequestTemplate),
			ResponseMappingTemplate: jsii.String(props.ResponseTemplate),
			PipelineConfig: &awsappsync.CfnResolver_PipelineConfigProperty{
				Functions: &functions,
			},
		},
	)
}
