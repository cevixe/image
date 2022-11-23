package iam

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func NewDynReadPol(tableArn string) awsiam.PolicyStatement {

	return awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: &[]*string{
			jsii.String("dynamodb:GetItem"),
			jsii.String("dynamodb:Scan"),
			jsii.String("dynamodb:Query"),
			jsii.String("dynamodb:BatchGetItem"),
			jsii.String("dynamodb:DescribeTable"),
		},
		Resources: &[]*string{
			jsii.String(tableArn),
		},
	})
}
