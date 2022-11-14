package appsync

import (
	"log"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"
)

type DSType uint8

const (
	DSType_None     DSType = 0
	DSType_Lambda   DSType = 1
	DSType_Dynamodb DSType = 2
)

type DataSourceProps struct {
	ApiId       string `field:"required"`
	Type        DSType `field:"required"`
	RoleArn     string `field:"optional"`
	LambdaArn   string `field:"optional"`
	TableRegion string `field:"optional"`
	TableName   string `field:"optional"`
	TableArn    string `field:"optional"`
}

func NewDataSource(mod module.Module, alias string, props *DataSourceProps) awsappsync.CfnDataSource {

	switch props.Type {
	case DSType_None:
		return newNoneDS(mod, alias, props)
	case DSType_Lambda:
		return newLambdaDS(mod, alias, props)
	case DSType_Dynamodb:
		return newDynamodbDS(mod, alias, props)
	default:
		log.Fatal("unknown appsync datasource type")
	}
	return nil
}

func newNoneDS(mod module.Module, alias string, props *DataSourceProps) awsappsync.CfnDataSource {

	name := naming.NewName(mod, naming.ResType_GraphQLDataSource, alias)

	return awsappsync.NewCfnDataSource(
		mod.Resource(),
		name.Logical(),
		&awsappsync.CfnDataSourceProps{
			Name:  jsii.String(alias),
			ApiId: jsii.String(props.ApiId),
			Type:  jsii.String("NONE"),
		},
	)
}

func newLambdaDS(mod module.Module, alias string, props *DataSourceProps) awsappsync.CfnDataSource {

	name := naming.NewName(mod, naming.ResType_GraphQLDataSource, alias)

	return awsappsync.NewCfnDataSource(
		mod.Resource(),
		name.Logical(),
		&awsappsync.CfnDataSourceProps{
			Name:  jsii.String(alias),
			ApiId: jsii.String(props.ApiId),
			Type:  jsii.String("AWS_LAMBDA"),
			LambdaConfig: &awsappsync.CfnDataSource_LambdaConfigProperty{
				LambdaFunctionArn: jsii.String(props.LambdaArn),
			},
			ServiceRoleArn: jsii.String(props.RoleArn),
		},
	)
}

func newDynamodbDS(mod module.Module, alias string, props *DataSourceProps) awsappsync.CfnDataSource {

	name := naming.NewName(mod, naming.ResType_GraphQLDataSource, alias)

	return awsappsync.NewCfnDataSource(
		mod.Resource(),
		name.Logical(),
		&awsappsync.CfnDataSourceProps{
			Name:           jsii.String(alias),
			ApiId:          jsii.String(props.ApiId),
			Type:           jsii.String("AMAZON_DYNAMODB"),
			ServiceRoleArn: jsii.String(props.RoleArn),
			DynamoDbConfig: awsappsync.CfnDataSource_DynamoDBConfigProperty{
				AwsRegion: jsii.String(props.TableRegion),
				TableName: jsii.String(props.TableName),
			},
		},
	)
}
