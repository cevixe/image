package function

import (
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewStateStoreUpdateFn(mod module.Module, props *FunctionProps) Function {

	name := export.StateStoreUpdateFn
	fn := appsync.NewFunction(mod, name, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  ssupdatefnrequest,
		ResponseTemplate: ssupdatefnresponse,
	})
	mod.Export(name, *fn.AttrFunctionId())
	return &impl{name: name, resource: fn}
}

const ssupdatefnrequest = `
#if( $ctx.stash.skip == true ) 
	#return($ctx.prev.result)
#end

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

#set( $updatedBy = $util.defaultIfNullOrBlank($ctx.identity.sub,"unknown") )
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

    ## Cevixe reserved properties
    #set( $reservedProps = ["__typename","id","version","createdAt","createdBy","updatedAt","updatedBy","__transaction","__status","__space"] )

    ## Increment "version" by 1 **
    $!{expAdd.put("version", ":one")}
    $!{expValues.put(":one", $util.dynamodb.toDynamoDB(1))}

    ## Iterate through each argument, skipping "id" and "expectedVersion" **
    #foreach( $entry in $util.map.copyAndRemoveAllKeys($args, $reservedProps).entrySet() )
        #if( $util.isNull($entry.value) )
            ## If the argument is set to "null", then remove that attribute from the item in DynamoDB **

            #set( $discard = ${expRemove.add("#${entry.key}")} )
            $!{expNames.put("#${entry.key}", "${entry.key}")}
        #else
            ## Otherwise set (or update) the attribute on the item in DynamoDB **

            $!{expSet.put("#${entry.key}", ":${entry.key}")}
            $!{expNames.put("#${entry.key}", "${entry.key}")}
            $!{expValues.put(":${entry.key}", $util.dynamodb.toDynamoDB($entry.value))}
        #end
    #end

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

    ## Start building the update expression, starting with attributes we're going to SET **
    #set( $expression = "" )
    #if( !${expSet.isEmpty()} )
        #set( $expression = "SET" )
        #foreach( $entry in $expSet.entrySet() )
            #set( $expression = "${expression} ${entry.key} = ${entry.value}" )
            #if ( $foreach.hasNext )
                #set( $expression = "${expression}," )
            #end
        #end
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
        "expression" : "#status = :expectedStatus",
        "expressionNames" : {
            "#status" : "__status"
        },
        "expressionValues" : {
            ":expectedStatus" : $util.dynamodb.toDynamoDBJson("alive")
        }
    }
    #else
    "condition" : {
        "expression" : "#status = :expectedStatus AND #version = :expectedVersion",
        "expressionNames" : {
            "#status" : "__status",
            "#version": "version"
        },
        "expressionValues" : {
            ":expectedStatus" : $util.dynamodb.toDynamoDBJson("alive"),
            ":expectedVersion" : $util.dynamodb.toDynamoDBJson($args.version)
        }
    }
    #end
}
`
const ssupdatefnresponse = `
#if($ctx.error)
    $util.error($ctx.error.message, $ctx.error.type)
#end
$util.toJson($ctx.result)
`
