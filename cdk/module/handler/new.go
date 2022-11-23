package handler

import (
	"fmt"
	"log"

	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/bus"
	"github.com/cevixe/cdk/module/function"
	"github.com/cevixe/cdk/service/iam"
	"github.com/cevixe/cdk/service/sns"
	"github.com/cevixe/cdk/service/sqs"
)

/*
Propiedades de configuración de una función handler.
*/
type HandlerProps struct {
	Name string `field:"required" json:"main"`
	/*
		Tipo de la función handler (HandlerType). Si no se especifica un valor
		se utilizará: HandlerType_Basic

		Ejemplos:
			-> HandlerType_Basic
			-> HandlerType_Standard
			-> HandlerType_Advanced
	*/
	Type HandlerType `field:"required" json:"type"`

	/*
		Tipos de eventos que deben ser escuchados y procesados por la función handler.
		Si se envía un array vacío, todos los eventos serán escuchados y procesados.
		Si no se envía un array, ningún evento será escuchado.

		Ejemplos:
			-> nil
			-> &[]string{}
			-> &[]string{"AccountCreated", "AccountUpdated", "AccountDeleted"}

	*/
	Events *[]string `field:"optional" json:"events"`

	/*
		Tipos de comandos que deben ser escuchados y procesados por la función handler.
		Si se envía un array vacío, todos los comandos serán escuchados y procesados.
		Si no se envía un array, ningún comando será escuchado.
		Ejemplos:
			-> nil
			-> &[]string{}
			-> &[]string{"CreateAccount", "UpdateAccount", "DeleteAccount"}

	*/
	Commands *[]string `field:"optional" json:"commands"`
}

func NewHandler(
	mod module.Module,
	advancedBus awssns.ITopic,
	standardBus awssns.ITopic,
	props *HandlerProps,
) Handler {

	entryFormat := "/cmd/handler/%s"
	entry := fmt.Sprintf(entryFormat, props.Name)
	fn := function.NewFunction(mod, props.Name, entry)

	commandStoreArn := mod.Import("core", export.CommandStoreArn)
	commandStoreName := mod.Import("core", export.CommandStoreName)

	stateStoreArn := mod.Import(mod.Name(), export.StateStoreArn)
	stateStoreName := mod.Import(mod.Name(), export.StateStoreName)

	objectStoreArn := mod.Import(mod.Name(), export.ObjectStoreArn)
	objectStoreName := mod.Import(mod.Name(), export.ObjectStoreName)

	fn.Resource().AddEnvironment(jsii.String("CVX_STATE_STORE"), jsii.String(stateStoreName), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_OBJECT_STORE"), jsii.String(objectStoreName), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_COMMAND_STORE"), jsii.String(commandStoreName), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_APP_NAME"), jsii.String(mod.App()), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_DOMAIN_NAME"), jsii.String(mod.Name()), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_HANDLER_NAME"), jsii.String(props.Name), nil)

	fn.Resource().AddToRolePolicy(iam.NewDynReadPol("*"))
	fn.Resource().AddToRolePolicy(iam.NewDynCrudPol(stateStoreArn))
	fn.Resource().AddToRolePolicy(iam.NewS3CrudPol(objectStoreArn))
	fn.Resource().AddToRolePolicy(iam.NewDynWritePol(commandStoreArn))

	filters := make([]*bus.Filter, 0)
	if props.Events != nil && len(*props.Events) > 0 {
		filters = append(filters, bus.NewFilter("event", *props.Events...))
	}
	if props.Commands != nil && len(*props.Commands) > 0 {
		filters = append(filters, bus.NewFilter("command", *props.Commands...))
	}
	if len(filters) == 0 {
		log.Fatalf("cannot determine message filter for handler: %s\n", props.Name)
	}

	switch props.Type {
	case HandlerType_Advanced:
		fn.Resource().AddEnvironment(jsii.String("CVX_HANDLER_MODE"), jsii.String("advanced"), nil)
		sns.NewSubscriptions(mod, props.Name, &sns.SubProps{
			Topic:    advancedBus,
			Function: fn.Resource(),
			Filters:  &filters,
			Queue:    sqs.NewQueue(mod, props.Name, sqs.QueueType_FIFO),
		})
	case HandlerType_Standard:
		fn.Resource().AddEnvironment(jsii.String("CVX_HANDLER_MODE"), jsii.String("advanced"), nil)
		sns.NewSubscriptions(mod, props.Name, &sns.SubProps{
			Topic:    standardBus,
			Function: fn.Resource(),
			Filters:  &filters,
			Queue:    sqs.NewQueue(mod, props.Name, sqs.QueueType_Standard),
		})
	case HandlerType_Basic:
		fn.Resource().AddEnvironment(jsii.String("CVX_HANDLER_MODE"), jsii.String("advanced"), nil)
		sns.NewSubscriptions(mod, props.Name, &sns.SubProps{
			Topic:    standardBus,
			Function: fn.Resource(),
			Filters:  &filters,
		})
	default:
		log.Fatalf("unsupported handler type: %v\n", props.Type)
	}

	return &handlerImpl{
		Function: fn,
		typ:      props.Type,
		events:   props.Events,
		commands: props.Commands,
	}
}
