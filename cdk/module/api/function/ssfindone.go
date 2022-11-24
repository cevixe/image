package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewStateStoreFindOneFn(mod module.Module, props *FunctionProps) Function {

	name := export.StateStoreFindOneFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  ssfindonefnrequest,
		ResponseTemplate: ssfindonefnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const ssfindonefnrequest = `
#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["__typename"]}) )
    $util.appendError("entity typename not specified", "EntityTypeNotFound")
	#return
#end

#if( $util.isNullOrBlank(${args["id"]}) )
    $util.appendError("entity id not specified", "EntityIdNotFound")
	#return
#end

{
    "version": "2018-05-29",
    "operation": "GetItem",
    "key": {
        "id": $util.dynamodb.toDynamoDBJson($ctx.args.id),
    }
}
`
const ssfindonefnresponse = `
#if($ctx.error)
    $util.appendError($ctx.error.message, $ctx.error.type)
	#return
#end
#if($ctx.result["__typename"] != $ctx.stash.input["__typename"])
	#return
#end
#if($ctx.result["__status"] == "dead")
	#return
#end
$util.toJson($ctx.result)
`
