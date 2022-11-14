package statestore

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/app/pkg/location"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/dynamodb"
	"github.com/cevixe/cdk/service/iam"
	"github.com/cevixe/cdk/service/lambda"
)

type StateStoreProps struct {
	AdvancedBusArn string
	StandardBusArn string
	Indexes        []string
}

func NewStateStore(mod module.Module, alias string, props *StateStoreProps) StateStore {

	indexes := make(map[string]*dynamodb.Key, 0)
	for _, idx := range props.Indexes {
		indexes[idx] = &dynamodb.Key{
			PartitionKey: dynamodb.NewAttribute(fmt.Sprintf("__%s-pk", idx), awsdynamodb.AttributeType_STRING),
			SortKey:      dynamodb.NewAttribute(fmt.Sprintf("__%s-sk", idx), awsdynamodb.AttributeType_STRING),
		}
	}
	indexes["by-section"] = &dynamodb.Key{
		PartitionKey: dynamodb.NewAttribute("__section", awsdynamodb.AttributeType_STRING),
		SortKey:      dynamodb.NewAttribute("id", awsdynamodb.AttributeType_STRING),
	}

	table := dynamodb.NewTable(mod, alias, &dynamodb.TableProps{
		Key: &dynamodb.Key{
			PartitionKey: dynamodb.NewAttribute("id", awsdynamodb.AttributeType_STRING),
		},
		GlobalIndexes: &indexes,
	})

	stream := awslambdaeventsources.NewDynamoEventSource(table, &awslambdaeventsources.DynamoEventSourceProps{
		Enabled:          jsii.Bool(true),
		BatchSize:        jsii.Number(10),
		StartingPosition: awslambda.StartingPosition_TRIM_HORIZON,
	})

	advancedcdc := lambda.NewGolangFunction(mod, "advcdc", location.EventCdcAdv)
	advancedcdc.AddEnvironment(jsii.String("CVX_EVENT_BUS"), &props.AdvancedBusArn, nil)
	advancedcdc.AddToRolePolicy(iam.NewSNSPublishPol(props.AdvancedBusArn))
	advancedcdc.AddEventSource(stream)

	standardcdc := lambda.NewGolangFunction(mod, "stdcdc", location.EventCdcStd)
	standardcdc.AddEnvironment(jsii.String("CVX_EVENT_BUS"), &props.StandardBusArn, nil)
	standardcdc.AddToRolePolicy(iam.NewSNSPublishPol(props.StandardBusArn))
	standardcdc.AddEventSource(stream)

	return &impl{
		module:      mod,
		name:        alias,
		resource:    table,
		advancedcdc: advancedcdc,
		standardcdc: standardcdc,
	}
}
