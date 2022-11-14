package core

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/bus"
	"github.com/cevixe/cdk/module/commandstore"
	"github.com/cevixe/cdk/module/eventstore"
)

func NewPlatform(scope constructs.Construct, app string) module.Module {

	mod := module.New(scope, module.Platform, app, "core")

	hostedZoneId := awscdk.NewCfnParameter(mod.Resource(), jsii.String("HostedZoneId"),
		&awscdk.CfnParameterProps{Type: jsii.String("String")})
	hostedZoneName := awscdk.NewCfnParameter(mod.Resource(), jsii.String("HostedZoneName"),
		&awscdk.CfnParameterProps{Type: jsii.String("String")})

	mod.Export(export.HostedZoneId, *hostedZoneId.ValueAsString())
	mod.Export(export.HostedZoneName, *hostedZoneName.ValueAsString())

	advancedbus := bus.NewBus(mod, "advancedbus",
		&bus.BusProps{
			Type: bus.BusType_Advanced,
		},
	)
	mod.Export(export.AdvancedBusName, *advancedbus.Resource().TopicName())
	mod.Export(export.AdvancedBusArn, *advancedbus.Resource().TopicArn())

	standardbus := bus.NewBus(mod, "standardbus",
		&bus.BusProps{
			Type: bus.BusType_Standard,
		},
	)
	mod.Export(export.StandardBusName, *standardbus.Resource().TopicName())
	mod.Export(export.StandardBusArn, *standardbus.Resource().TopicArn())

	commandstore := commandstore.NewCommandStore(mod, "commandstore",
		&commandstore.CommandStoreProps{
			AdvancedBus: advancedbus,
			StandardBus: standardbus,
		},
	)
	mod.Export(export.CommandStoreName, *commandstore.Resource().TableName())
	mod.Export(export.CommandStoreArn, *commandstore.Resource().TableArn())

	eventstore := eventstore.NewEventStore(mod, "eventstore",
		&eventstore.EventStoreProps{
			AdvancedBus: advancedbus,
		},
	)
	mod.Export(export.EventStoreName, *eventstore.Resource().TableName())
	mod.Export(export.EventStoreArn, *eventstore.Resource().TableArn())

	return mod
}
