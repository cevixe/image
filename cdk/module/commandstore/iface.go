package commandstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/cevixe/cdk/module"
)

type CommandStore interface {
	Module() module.Module
	Name() string
	Resource() awsdynamodb.Table
	AdvancedCdc() awslambda.Function
	StandardCdc() awslambda.Function
}
