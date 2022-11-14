package resolver

import "github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"

type Resolver interface {
	Name() string
	Resource() awsappsync.CfnResolver
}
