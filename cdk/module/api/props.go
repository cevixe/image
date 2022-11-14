package api

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/cevixe/cdk/module/objectstore"
	"github.com/cevixe/cdk/module/statestore"
)

type ApiProps struct {
	Zone        awsroute53.IHostedZone  `field:"required"`
	StateStore  statestore.StateStore   `field:"required"`
	ObjectStore objectstore.ObjectStore `field:"required"`
}

type DataSourceType string

const (
	DataSourceType_Function = "function"
	DataSourceType_Store    = "store"
	DataSourceType_Mock     = "mock"
)

type ApiProps_DataSource struct {
	Name string         `field:"required"`
	Type DataSourceType `field:"required"`
}

type ApiProps_Function struct {
	Name       string `field:"required"`
	DataSource string `field:"required"`
}

type ApiProps_Resolver struct {
	Name      string   `field:"required"`
	Operation string   `field:"required"`
	Functions []string `field:"required"`
}

type ApiConfigProps struct {
	DataSources []ApiProps_DataSource `field:"required"`
	Functions   []ApiProps_Function   `field:"required"`
	Resolvers   []ApiProps_Resolver   `field:"required"`
}
