package datasource

import "github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"

type impl struct {
	name     string
	resource awsappsync.CfnDataSource
}

func (d *impl) Name() string {
	return d.name
}

func (d *impl) Resource() awsappsync.CfnDataSource {
	return d.resource
}
