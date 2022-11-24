package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewObjectStoreUploadFn(mod module.Module, props *FunctionProps) Function {

	name := export.ObjectStoreUploadFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  osuploadfnrequest,
		ResponseTemplate: osuploadfnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const osuploadfnrequest = `
#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["space"]}) )
    $util.appendError("object space not specified", "ObjectSpaceNotFound")
	#return
#end
#set( $space = ${args["space"]} )

#set( $timezone = "America/Lima" )
#if( !$util.isNullOrBlank(${args["timezone"]}) )
	#set( $timezone = ${args["timezone"]} )
#end

#set( $id = $util.autoUlid() )
#set( $now = $util.time.nowEpochMilliSeconds() )
#set( $author = $util.defaultIfNullOrBlank($ctx.identity.username,"unknown") )

#set( $year = $util.time.epochMilliSecondsToFormatted($now, "yyyy", $timezone) )
#set( $month = $util.time.epochMilliSecondsToFormatted($now, "MM", $timezone) )
#set( $day = $util.time.epochMilliSecondsToFormatted($now, "dd", $timezone) )

#set( $directory = "$space/$year/$month/$day" )
#set( $filename = "$directory/$util.autoUlid()" )
#if( !$util.isNullOrBlank(${args["extension"]}) )
	#set( $extension = ${args["extension"]} )
	#set( $filename = "$filename.$extension" )
#end

{
  "version" : "2017-02-28",
  "operation": "Invoke",
  "payload": $util.toJson({
    "operation": "upload",
    "name": "$filename"
  })
}
`

const osuploadfnresponse = `
#if($ctx.error)
  $util.appendError($ctx.error.message, $ctx.error.type)
  #return
#end
$util.toJson($context.result)
`
