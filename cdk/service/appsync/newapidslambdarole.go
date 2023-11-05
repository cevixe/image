package appsync

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/iam"
)

type ApiDsRoleProps struct {
	Alias string `field:"required"`
	Arn   string `field:"required"`
}

func NewApiDsLambdaRole(mod module.Module, apiDsRoleProps *ApiDsRoleProps) awsiam.Role {
	alias := fmt.Sprintf("%s_%s_%s", mod.Name(), "gqapidslamba", apiDsRoleProps.Alias)
	role := iam.NewServiceRole(mod, alias, "appsync.amazonaws.com")
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: &[]*string{
			jsii.String("lambda:InvokeFunction"),
		},
		Resources: &[]*string{
			jsii.String(apiDsRoleProps.Arn),
		},
	}))
	return role
}
