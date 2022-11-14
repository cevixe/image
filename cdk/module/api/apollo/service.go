package apollo

import (
	"fmt"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

type ServiceResolverProps struct {
	ApiId  string `field:"required"`
	Schema string `field:"required"`
}

const ServiceRequest = `
{
    "version": "2017-02-28",
	"payload": {
		"sdl": "%s"
	}
}
`

const ServiceResponse = `
$util.toJson($context.result)
`

func NewServiceResolver(mod module.Module, props *ServiceResolverProps) awsappsync.CfnResolver {

	escapedSchema := escapeSchema(props.Schema)

	request := fmt.Sprintf(ServiceRequest, *escapedSchema)
	response := ServiceResponse

	fn := appsync.NewFunction(mod, "apolloservicefn", &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   "mock",
		RequestTemplate:  request,
		ResponseTemplate: response,
	})

	rs := appsync.NewResolver(mod, "apolloservicers", &appsync.ResolverProps{
		ApiId:            props.ApiId,
		Type:             "Query",
		Field:            "_service",
		RequestTemplate:  `{}`,
		ResponseTemplate: `$ctx.result`,
		Functions:        []string{*fn.AttrFunctionId()},
	})
	rs.AddDependsOn(fn)

	return rs
}

func escapeSchema(schema string) *string {
	escapedSchema := schema
	escapedSchema = strings.ReplaceAll(escapedSchema, `"`, `\"`)
	escapedSchema = strings.ReplaceAll(escapedSchema, "\n", `\n`)
	return &escapedSchema
}
