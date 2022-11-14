package domain

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/module/api"
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

	NewService(scope, &ServiceProps{
		App:  props.App.Name,
		Name: name,
		Api: api.ApiConfigProps{
			Resolvers: resolvers,
			Functions: functions,
		},
	})
}
