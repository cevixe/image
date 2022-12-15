package apollo

import (
	"fmt"

	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
)

type EntitiesResolverProps struct {
	ApiId          string `field:"required"`
	StateStoreName string `field:"required"`
}

const EntitiesRequest = `
#set($ids = [])
#foreach($item in ${ctx.args.representations})
	#set($map = {})
	$util.qr($map.put("id", $util.dynamodb.toString($item.id)))
	$util.qr($ids.add($map))
#end

{
	"version" : "2018-05-29",
	"operation" : "BatchGetItem",
	"tables" : {
    	"%s": {
    		"keys": $util.toJson($ids),
			"consistentRead": true
		}
    }
}
`

const EntitiesResponse = `
#if($ctx.error)
	$util.error($ctx.error.message, $ctx.error.type)
#end
#set($items = $context.result.data.%s)
$util.toJson($items)
`

func NewEntitiesResolver(mod module.Module, props *EntitiesResolverProps) awsappsync.CfnResolver {

	request := fmt.Sprintf(EntitiesRequest, props.StateStoreName)
	response := fmt.Sprintf(EntitiesResponse, props.StateStoreName)

	fn := appsync.NewFunction(mod, "apolloentitiesfn", &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   "statestore",
		RequestTemplate:  request,
		ResponseTemplate: response,
	})

	rs := appsync.NewResolver(mod, "apolloentitiesrs", &appsync.ResolverProps{
		ApiId:            props.ApiId,
		Type:             "Query",
		Field:            "_entities",
		RequestTemplate:  `{}`,
		ResponseTemplate: `$util.toJson($ctx.result)`,
		Functions:        []string{*fn.AttrFunctionId()},
	})
	rs.AddDependsOn(fn)

	return rs
}
