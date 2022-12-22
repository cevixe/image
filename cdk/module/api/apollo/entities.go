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
#set($set = {})
#set($ids = [])
#foreach($item in ${ctx.args.representations})
	#if( !$set.containsKey($item.id) )
      $util.qr($set.put($item.id, $item.id))
      #set($map = {})
      $util.qr($map.put("id", $util.dynamodb.toString($item.id)))
      $util.qr($ids.add($map))
    #end
#end


{
	"version" : "2018-05-29",
	"operation" : "BatchGetItem",
	"tables" : {
    	"%s": {
    		"keys": $util.toJson($ids),
			"consistentRead": false
		}
    }
}
`

const EntitiesResponse = `
#if($ctx.error)
	$util.error($ctx.error.message, $ctx.error.type)
#end
#set($items = $context.result.data.%s)

#set($map = {})
#foreach($item in $items)
	$util.qr($map.put($item.id, $item))
#end

#set($result = [])
#foreach($item in ${ctx.args.representations})
	$util.qr($result.add($map.get($item.id)))
#end

$util.toJson($result)
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
