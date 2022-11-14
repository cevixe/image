package datasource

import "github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"

type DataSource interface {
	Name() string
	Resource() awsappsync.CfnDataSource
}
