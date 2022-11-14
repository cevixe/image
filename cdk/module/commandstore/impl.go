package commandstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/cevixe/cdk/module"
)

type impl struct {
	module      module.Module
	name        string
	resource    awsdynamodb.Table
	advancedcdc awslambda.Function
	standardcdc awslambda.Function
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

func (store *impl) AdvancedCdc() awslambda.Function {
	return store.advancedcdc
}

func (store *impl) StandardCdc() awslambda.Function {
	return store.standardcdc
}
