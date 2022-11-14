package iam

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func NewSNSPublishPol(topicArn string) awsiam.PolicyStatement {

	return awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: &[]*string{
			jsii.String("sns:Publish"),
		},
		Resources: &[]*string{
			jsii.String(topicArn),
		},
	})
}
