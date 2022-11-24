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
#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["__typename"]}) )
    $util.appendError("entity typename not specified", "EntityTypeNotFound")
    #return
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
    $util.appendError($ctx.error.message, $ctx.error.type)
    #return
#end
$util.toJson($ctx.result)
`
