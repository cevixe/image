package datasource

import "github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"

type DSType uint8

const (
	DSType_Mock   DSType = 0
	DSType_Table  DSType = 1
	DSType_Lambda DSType = 2
)

type DataSourceProps struct {
	ApiId       string             `field:"required"`
	Type        DSType             `field:"required"`
	RoleArn     string             `field:"optional"`
	Region      string             `field:"optional"`
	StateStore  awsdynamodb.ITable `field:"optional"`
	FunctionArn string             `field:"optional"`
}
