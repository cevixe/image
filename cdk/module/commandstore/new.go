package commandstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/app/pkg/location"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/bus"
	"github.com/cevixe/cdk/service/dynamodb"
	"github.com/cevixe/cdk/service/iam"
	"github.com/cevixe/cdk/service/lambda"
)

type CommandStoreProps struct {
	AdvancedBus bus.Bus
	StandardBus bus.Bus
}

func NewCommandStore(mod module.Module, alias string, props *CommandStoreProps) CommandStore {

	table := dynamodb.NewTable(mod, alias, &dynamodb.TableProps{
		Key: &dynamodb.Key{
			PartitionKey: dynamodb.NewAttribute("id", awsdynamodb.AttributeType_STRING),
		},
	})

	stream := awslambdaeventsources.NewDynamoEventSource(table, &awslambdaeventsources.DynamoEventSourceProps{
		Enabled:          jsii.Bool(true),
		BatchSize:        jsii.Number(10),
		StartingPosition: awslambda.StartingPosition_TRIM_HORIZON,
	})

	advancedbus := props.AdvancedBus.Resource()
	advancedcdc := lambda.NewGolangFunction(mod, "advcdc", location.CommandCdcAdv)
	advancedcdc.AddEnvironment(jsii.String("CVX_EVENT_BUS"), advancedbus.TopicArn(), nil)
	advancedcdc.AddToRolePolicy(iam.NewSNSPublishPol(*advancedbus.TopicArn()))
	advancedcdc.AddEventSource(stream)

	standardbus := props.StandardBus.Resource()
	standardcdc := lambda.NewGolangFunction(mod, "stdcdc", location.CommandCdcStd)
	standardcdc.AddEnvironment(jsii.String("CVX_EVENT_BUS"), standardbus.TopicArn(), nil)
	standardcdc.AddToRolePolicy(iam.NewSNSPublishPol(*standardbus.TopicArn()))
	standardcdc.AddEventSource(stream)

	return &impl{
		module:      mod,
		name:        alias,
		resource:    table,
		advancedcdc: advancedcdc,
		standardcdc: standardcdc,
	}
}
