package eventstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/app/pkg/location"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/bus"
	"github.com/cevixe/cdk/service/dynamodb"
	"github.com/cevixe/cdk/service/iam"
	"github.com/cevixe/cdk/service/lambda"
	"github.com/cevixe/cdk/service/sns"
	"github.com/cevixe/cdk/service/sqs"
)

type EventStoreProps struct {
	AdvancedBus bus.Bus
}

func NewEventStore(mod module.Module, alias string, props *EventStoreProps) EventStore {

	table := dynamodb.NewTable(mod, alias, &dynamodb.TableProps{
		Key: &dynamodb.Key{
			PartitionKey: dynamodb.NewAttribute("source", awsdynamodb.AttributeType_STRING),
			SortKey:      dynamodb.NewAttribute("id", awsdynamodb.AttributeType_STRING),
		},
		GlobalIndexes: &map[string]*dynamodb.Key{
			"by-type": {
				PartitionKey: dynamodb.NewAttribute("type", awsdynamodb.AttributeType_STRING),
				SortKey:      dynamodb.NewAttribute("time", awsdynamodb.AttributeType_STRING),
			},
			"by-user": {
				PartitionKey: dynamodb.NewAttribute("iocevixeuser", awsdynamodb.AttributeType_STRING),
				SortKey:      dynamodb.NewAttribute("time", awsdynamodb.AttributeType_STRING),
			},
			"by-transaction": {
				PartitionKey: dynamodb.NewAttribute("iocevixetransaction", awsdynamodb.AttributeType_STRING),
				SortKey:      dynamodb.NewAttribute("time", awsdynamodb.AttributeType_STRING),
			},
		},
	})

	eventhandler := lambda.NewGolangFunction(mod, "eventhdl", location.EventHandler)
	eventhandler.AddEnvironment(jsii.String("CVX_EVENT_STORE"), table.TableName(), nil)
	eventhandler.AddToRolePolicy(iam.NewDynWritePol(*table.TableArn()))

	sns.NewSubscriptions(mod, "eventhdl", &sns.SubProps{
		Topic:    props.AdvancedBus.Resource(),
		Function: eventhandler,
		Filters:  &[]*bus.Filter{bus.NewFilter("event")},
		Queue:    sqs.NewQueue(mod, "eventhdl", sqs.QueueType_FIFO),
	})

	return &impl{
		module:   mod,
		name:     alias,
		resource: table,
		handler:  eventhandler,
	}
}
