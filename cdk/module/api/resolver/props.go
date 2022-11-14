package resolver

import "github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"

type ResolverProps struct {
	ApiId     string                      `field:"required"`
	Type      string                      `field:"required"`
	Field     string                      `field:"required"`
	Functions []string                    `field:"required"`
	Schema    awsappsync.CfnGraphQLSchema `field:"required"`
}
