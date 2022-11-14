package function

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
)

type impl struct {
	name     string
	resource awsappsync.CfnFunctionConfiguration
}

func (a *impl) Name() string {
	return a.name
}

func (a *impl) Resource() awsappsync.CfnFunctionConfiguration {
	return a.resource
}
