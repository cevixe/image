package function

import "github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"

type Function interface {
	Name() string
	Resource() awsappsync.CfnFunctionConfiguration
}
