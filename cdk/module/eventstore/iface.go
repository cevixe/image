package eventstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/cevixe/cdk/module"
)

type EventStore interface {
	Module() module.Module
	Name() string
	Resource() awsdynamodb.Table
	Handler() awslambda.Function
}
