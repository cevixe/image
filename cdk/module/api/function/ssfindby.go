package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewStateStoreFindByFn(mod module.Module, props *FunctionProps) Function {

	name := export.StateStoreFindByFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  ssfindbyfnrequest,
		ResponseTemplate: ssfindbyfnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const ssfindbyfnrequest = `
#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["__typename"]}) )
    $util.appendError("entity typename not specified", "EntityTypeNotFound")
    #return
#end

#if( $util.isNullOrBlank(${args["indexName"]}) )
    $util.appendError("entity index name not specified", "EntityIndexNameNotFound")
    #return
#end

#if( $util.isNullOrBlank(${args["indexValue"]}) )
    $util.appendError("entity index value not specified", "EntityIndexValueNotFound")
    #return
#end

#set( $typename = ${args["__typename"]} )
#set( $indexName =  ${args["indexName"]} )
#set( $indexPrefix =  "__" )
#set( $indexSuffix =  "-pk" )
#set( $indexPkName=  "$indexPrefix$indexName$indexSuffix" )
#set( $indexPkValue =  ${args["indexValue"]} )

{
    "version": "2018-05-29",
    "operation" : "Query",
    "index" : "$indexName",
    "query" : {
        "expression": "#indexPk = :indexPk",
        "expressionNames" : {
            "#indexPk" : "$indexPkName"
        },
         "expressionValues" : {
            ":indexPk" : $util.dynamodb.toDynamoDBJson($indexPkValue)
        }
    },
    "scanIndexForward": false,
    "limit": $util.defaultIfNull(${args.limit}, 20),
    "nextToken": $util.toJson($util.defaultIfNullOrBlank($args.nextToken, null))
}
`
const ssfindbyfnresponse = `
#if($ctx.error)
    $util.appendError($ctx.error.message, $ctx.error.type)
    #return
#end
$util.toJson($ctx.result)
`
