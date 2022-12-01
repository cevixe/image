package api

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/cevixe/cdk/module"
)

type Api interface {
	Module() module.Module
	Name() string
	URL() string
	Role() awsiam.Role
	Schema() awsappsync.CfnGraphQLSchema
	Resource() awsappsync.CfnGraphQLApi
}
