package lambda

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/iam"
)

func NewFunctionRole(mod module.Module, alias string) awsiam.Role {
	role := iam.NewServiceRole(mod, alias, "lambda.amazonaws.com")
	role.AddToPrincipalPolicy(awsiam.NewPolicyStatement(
		&awsiam.PolicyStatementProps{
			Effect: awsiam.Effect_ALLOW,
			Actions: &[]*string{
				jsii.String("logs:CreateLogGroup"),
				jsii.String("logs:CreateLogStream"),
				jsii.String("logs:PutLogEvents"),
			},
			Resources: &[]*string{
				jsii.String("*"),
			},
		}))
	return role
}
