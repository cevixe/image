package resolver

import "github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"

type impl struct {
	name     string
	resource awsappsync.CfnResolver
}

func (r *impl) Name() string {
	return r.name
}

func (r *impl) Resource() awsappsync.CfnResolver {
	return r.resource
}
