package eventstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/cevixe/cdk/module"
)

type impl struct {
	module   module.Module
	name     string
	resource awsdynamodb.Table
	handler  awslambda.Function
}

func (store *impl) Module() module.Module {
	return store.module
}

func (store *impl) Name() string {
	return store.name
}

func (store *impl) Resource() awsdynamodb.Table {
	return store.resource
}

func (store *impl) Handler() awslambda.Function {
	return store.handler
}
