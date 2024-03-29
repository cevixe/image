package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewStateStoreCreateFn(mod module.Module, props *FunctionProps) Function {

	name := export.StateStoreCreateFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  sscreatefnrequest,
		ResponseTemplate: sscreatefnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const sscreatefnrequest = `
#if( $ctx.stash.skip == true ) 
	#return($ctx.prev.result)
#end

#set( $args = $ctx.stash.input )

#if($util.isNullOrBlank(${args["__typename"]}))
    $util.error("entity typename not specified", "EntityTypeNotFound")
#end

#set( $typename = ${args["__typename"]} )
#set( $id = $util.autoUlid() )
#if( !$util.isNullOrBlank(${args["id"]}) )
	#set( $id = ${args["id"]} )
#end

#set( $createdBy = $util.defaultIfNullOrBlank($ctx.identity.sub,"unknown") )
#set( $createdAt = $util.time.nowISO8601() )
#set( $space = "alive#$typename" )

#set( $tracingHeader = $context.request.headers["x-amzn-trace-id"] )
#set( $transaction = $util.defaultIfNullOrBlank($tracingHeader.replaceAll("Root=", ""), $util.autoId()) )

$util.qr( $args.put("version", 1) )
$util.qr( $args.put("createdBy", $createdBy) )
$util.qr( $args.put("createdAt", $createdAt) )
$util.qr( $args.put("updatedBy", $createdBy) )
$util.qr( $args.put("updatedAt", $createdAt) )

$util.qr( $args.put("__typename", $typename) )
$util.qr( $args.put("__transaction", $transaction) )
$util.qr( $args.put("__status", "alive") )
$util.qr( $args.put("__space", $space) )

#set( $attributes = $util.dynamodb.toMapValues($args) )

{
    "version" : "2018-05-29",
    "operation" : "PutItem",
    "key" : {
        "id" : { "S" : "${id}" }
    },
    "attributeValues": $util.toJson($attributes),
	"condition" : {
        "expression" : "attribute_not_exists(id)"
    }
}
`
const sscreatefnresponse = `
#if($ctx.error)
    $util.error($ctx.error.message, $ctx.error.type)
#end
#set($result = $ctx.result)
$util.qr($result.put("updatedBy", { "__typename": "User", "id": "$result.updatedBy" }))
$util.qr($result.put("createdBy", { "__typename": "User", "id": "$result.createdBy" }))
$util.toJson($result)
`
