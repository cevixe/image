package iam

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func NewDynWritePol(tableArn string) awsiam.PolicyStatement {

	return awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: &[]*string{
			jsii.String("dynamodb:PutItem"),
			jsii.String("dynamodb:UpdateItem"),
			jsii.String("dynamodb:BatchWriteItem"),
			jsii.String("dynamodb:ConditionCheckItem"),
		},
		Resources: &[]*string{
			jsii.String(tableArn),
		},
	})
}
