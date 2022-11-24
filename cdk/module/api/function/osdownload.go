package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewObjectStoreDownloadFn(mod module.Module, props *FunctionProps) Function {

	name := export.ObjectStoreDownloadFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  osdownloadfnrequest,
		ResponseTemplate: osdownloadfnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const osdownloadfnrequest = `

#if( $ctx.stash.skip == true ) 
	#return($ctx.prev.result)
#end

#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["name"]}) )
    $util.error("object name not specified", "ObjectNameNotFound")
#end

#set( $name = ${args["name"]} )

{  
	"version" : "2017-02-28",
	"operation": "Invoke",
	"payload": $util.toJson({
	  "operation": "download",
	  "name": "$name"
	})
  }
`

const osdownloadfnresponse = `
#if($ctx.error)
  $util.error($ctx.error.message, $ctx.error.type)
#end
$util.toJson($context.result)
`
