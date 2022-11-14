package api

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/cevixe/cdk/module"
)

type apiImpl struct {
	module   module.Module
	name     string
	record   string
	domain   string
	key      awsappsync.CfnApiKey
	schema   awsappsync.CfnGraphQLSchema
	role     awsiam.Role
	resource awsappsync.CfnGraphQLApi
}

func (a *apiImpl) Module() module.Module {
	return a.module
}

func (a *apiImpl) Name() string {
	return a.name
}

func (a *apiImpl) Schema() awsappsync.CfnGraphQLSchema {
	return a.schema
}

func (a *apiImpl) Key() awsappsync.CfnApiKey {
	return a.key
}

func (a *apiImpl) Role() awsiam.Role {
	return a.role
}

func (a *apiImpl) URL() string {
	return fmt.Sprintf("https://%s.%s/graphql", a.record, a.domain)
}

func (a *apiImpl) Resource() awsappsync.CfnGraphQLApi {
	return a.resource
}
