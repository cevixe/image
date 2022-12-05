package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewStateStoreFindAllFn(mod module.Module, props *FunctionProps) Function {

	name := export.StateStoreFindAllFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  ssfindallfnrequest,
		ResponseTemplate: ssfindallfnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const ssfindallfnrequest = `
#if( $ctx.stash.skip == true ) 
	#return($ctx.prev.result)
#end

#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["__typename"]}) )
    $util.error("entity typename not specified", "EntityTypeNotFound")
#end

#set( $typename = ${args["__typename"]} )
#set( $space = "alive#$typename" )

{
    "version": "2018-05-29",
    "operation" : "Query",
    "index" : "by-space",
    "query" : {
        "expression": "#space = :space",
        "expressionNames" : {
            "#space" : "__space"
        },
         "expressionValues" : {
            ":space" : $util.dynamodb.toDynamoDBJson($space)
        }
    },
    "scanIndexForward": false,
    "limit": $util.defaultIfNull(${args.limit}, 20),
    "nextToken": $util.toJson($util.defaultIfNullOrBlank($args.nextToken, null))
}
`
const ssfindallfnresponse = `
#if($ctx.error)
    $util.error($ctx.error.message, $ctx.error.type)
#end

#set($result = [])
#foreach( $item in $ctx.result )
    $!{item.put("updatedBy", { "__typename": "User", "id": "$item.updatedBy" })}
    $!{item.put("createdBy", { "__typename": "User", "id": "$item.createdBy" })}
    #set( $discard = ${result.add($item)} )
#end
$util.toJson($result)
`
