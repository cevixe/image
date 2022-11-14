package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewStateStoreDeleteFn(mod module.Module, props *FunctionProps) Function {

	name := export.StateStoreDeleteFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  ssdeletefnrequest,
		ResponseTemplate: ssdeletefnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const ssdeletefnrequest = `
#set( $args = $ctx.stash.input )

#if( $util.isNullOrBlank(${args["__typename"]}) )
    $util.error("entity typename not specified", "EntityTypeNotFound")
#end

#if( $util.isNullOrBlank(${args["id"]}) )
    $util.error("entity id not specified", "EntityIdNotFound")
#end

#if( !$util.isNullOrBlank(${args["version"]}) && !$util.isNumber(${args["version"]}) )
    $util.error("entity version not numeric", "EntityVersionNotNumeric")
#end

#set( $typename = ${args["__typename"]} )
#set( $section = "archived#${typename}" )

#set( $updatedBy = $util.defaultIfNullOrBlank($ctx.identity.username,"unknown") )
#set( $updatedAt = $util.time.nowISO8601() )

#set( $tracingHeader = $context.request.headers["x-amzn-trace-id"] )
#set( $transaction = $util.defaultIfNullOrBlank($tracingHeader.replaceAll("Root=", ""), $util.autoId()) )

{
    "version" : "2018-05-29",
    "operation" : "UpdateItem",
    "key" : {
        "id" : $util.dynamodb.toDynamoDBJson($args.id)
    },

    ## Set up some space to keep track of things we're updating **
    #set( $expNames  = {} )
    #set( $expValues = {} )
    #set( $expSet = {} )
    #set( $expAdd = {} )
    #set( $expRemove = [] )

    ## Increment "version" by 1 **

    $!{expAdd.put("version", ":one")}
    $!{expValues.put(":one", $util.dynamodb.toDynamoDB(1))}

    ## Set cevixe reserved properties

    $!{expSet.put("#updatedAt", ":updatedAt")}
    $!{expNames.put("#updatedAt", "updatedAt")}
    $!{expValues.put(":updatedAt", $util.dynamodb.toDynamoDB($updatedAt))}

    $!{expSet.put("#updatedBy", ":updatedBy")}
    $!{expNames.put("#updatedBy", "updatedBy")}
    $!{expValues.put(":updatedBy", $util.dynamodb.toDynamoDB($updatedBy))}

    $!{expSet.put("#__transaction", ":__transaction")}
    $!{expNames.put("#__transaction", "__transaction")}
    $!{expValues.put(":__transaction", $util.dynamodb.toDynamoDB($transaction))}

    $!{expSet.put("#__archived", ":__archived")}
    $!{expNames.put("#__archived", "__archived")}
    $!{expValues.put(":__archived", $util.dynamodb.toDynamoDB(true))}

    $!{expSet.put("#__section", ":__section")}
    $!{expNames.put("#__section", "__section")}
    $!{expValues.put(":__section", $util.dynamodb.toDynamoDB($section))}

    ## Cevixe reserved properties
    #set( $reservedIdx = ["by-section"] )

    ## Iterate through indexes
    #foreach( $idx in $util.map.copyAndRemoveAll($args.indexes, $reservedIdx) )
        #set( $pk = "__${idx}-pk" )
        #set( $sk = "__${idx}-sk" )
        #set( $discard = ${expRemove.add("#${pk}")} )
        #set( $discard = ${expRemove.add("#${sk}")} )
        $!{expNames.put("#${pk}", "${pk}")}
        $!{expNames.put("#${sk}", "${sk}")}
    #end

    ## Continue building the update expression, adding attributes we're going to ADD **
    #if( !${expAdd.isEmpty()} )
        #set( $expression = "${expression} ADD" )
        #foreach( $entry in $expAdd.entrySet() )
            #set( $expression = "${expression} ${entry.key} ${entry.value}" )
            #if ( $foreach.hasNext )
                #set( $expression = "${expression}," )
            #end
        #end
    #end

    ## Continue building the update expression, adding attributes we're going to REMOVE **
    #if( !${expRemove.isEmpty()} )
        #set( $expression = "${expression} REMOVE" )

        #foreach( $entry in $expRemove )
            #set( $expression = "${expression} ${entry}" )
            #if ( $foreach.hasNext )
                #set( $expression = "${expression}," )
            #end
        #end
    #end

    ## Finally, write the update expression into the document, along with any expressionNames and expressionValues **
    "update" : {
        "expression" : "${expression}",
        #if( !${expNames.isEmpty()} )
            "expressionNames" : $utils.toJson($expNames),
        #end
        #if( !${expValues.isEmpty()} )
            "expressionValues" : $utils.toJson($expValues),
        #end
    },
    #if( $util.isNullOrBlank(${args.version}) )
    "condition" : {
        "expression" : "#archived = :expectedArchived",
        "expressionNames" : {
            "#archived" : "__archived"
        },
        "expressionValues" : {
            ":expectedArchived" : $util.dynamodb.toDynamoDBJson(false)
        }
    }
    #else
    "condition" : {
        "expression" : "#archived = :expectedArchived AND #version = :expectedVersion",
        "expressionNames" : {
            "#archived" : "__archived",
            "#version": "version"
        },
        "expressionValues" : {
            ":expectedArchived" : $util.dynamodb.toDynamoDBJson(false),
            ":expectedVersion" : $util.dynamodb.toDynamoDBJson($args.version)
        }
    }
    #end
}
`
const ssdeletefnresponse = `
#if($ctx.error)
    $util.error($ctx.error.message, $ctx.error.type)
#end
$util.toJson($ctx.result)
`
