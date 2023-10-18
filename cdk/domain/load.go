package domain

import (
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/module/api"
	"github.com/cevixe/cdk/module/handler"
	"github.com/cevixe/cdk/spec/v20221023"
)

func Load20221023(scope constructs.Construct, name string, props *spec.Properties) {

	functions := make([]api.ApiProps_Function, 0)
	for _, item := range props.Api.Functions {
		functions = append(functions, api.ApiProps_Function{
			Name:       item.Name,
			DataSource: item.DataSource,
		})
	}

	resolvers := make([]api.ApiProps_Resolver, 0)
	for _, item := range props.Api.Resolvers {
		resolvers = append(resolvers, api.ApiProps_Resolver{
			Name:      item.Name,
			Operation: item.Operation,
			Functions: item.Functions,
		})
	}

	datasources := make([]api.ApiProps_DataSource, 0)
	for _, item := range props.Api.DataSources {
		datasources = append(datasources, api.ApiProps_DataSource{
			Name: item.Name,
			Type: api.DataSourceType(item.Type),
		})
	}

	handlers := make([]handler.HandlerProps, 0)
	for idx := range props.Handlers {
		item := props.Handlers[idx]
		var handlerType handler.HandlerType
		switch item.Type {
		case spec.HandlerType_Basic:
			handlerType = handler.HandlerType_Basic
		case spec.HandlerType_Standard:
			handlerType = handler.HandlerType_Standard
		case spec.HandlerType_Advanced:
			handlerType = handler.HandlerType_Advanced
		default:
			log.Fatalf("unsupport handler type: %s\n", item.Type)
		}
		handlers = append(handlers, handler.HandlerProps{
			Name:     item.Name,
			Type:     handlerType,
			Events:   &item.Events,
			Commands: &item.Commands,
		})
	}

	NewService(scope, &ServiceProps{
		App:  props.App.Name,
		Name: name,
		Api: api.ApiConfigProps{
			Resolvers:   resolvers,
			Functions:   functions,
			DataSources: datasources,
		},
		Handlers: handlers,
	})
}
