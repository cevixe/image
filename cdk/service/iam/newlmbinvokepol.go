package iam

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func NewLambdaInvokePol(lambdaArn string) awsiam.PolicyStatement {

	return awsiam.NewPolicyStatement(
		&awsiam.PolicyStatementProps{
			Effect: awsiam.Effect_ALLOW,
			Actions: &[]*string{
				jsii.String("lambda:invokeFunction"),
			},
			Resources: &[]*string{
				jsii.String(lambdaArn),
			},
		},
	)
}
