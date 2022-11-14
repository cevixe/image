package datasource

import (
	"log"

	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewDataSource(mod module.Module, alias string, props *DataSourceProps) DataSource {
	switch props.Type {
	case DSType_Mock:
		return newMockDataSource(mod, alias, props)
	case DSType_Table:
		return newTableDataSource(mod, alias, props)
	case DSType_Lambda:
		return newLambdaDataSource(mod, alias, props)
	default:
		log.Fatalf("api datasource type not supported")
		return nil
	}
}

func newMockDataSource(mod module.Module, alias string, props *DataSourceProps) DataSource {
	ds := appsync.NewDataSource(mod, alias, &appsync.DataSourceProps{
		ApiId: props.ApiId,
		Type:  appsync.DSType_None,
	})
	return &impl{name: alias, resource: ds}
}

func newTableDataSource(mod module.Module, alias string, props *DataSourceProps) DataSource {
	ds := appsync.NewDataSource(mod, alias, &appsync.DataSourceProps{
		ApiId:       props.ApiId,
		RoleArn:     props.RoleArn,
		Type:        appsync.DSType_Dynamodb,
		TableRegion: props.Region,
		TableName:   *props.StateStore.TableName(),
		TableArn:    *props.StateStore.TableArn(),
	})
	return &impl{name: alias, resource: ds}
}

func newLambdaDataSource(mod module.Module, alias string, props *DataSourceProps) DataSource {
	ds := appsync.NewDataSource(mod, alias, &appsync.DataSourceProps{
		ApiId:     props.ApiId,
		RoleArn:   props.RoleArn,
		Type:      appsync.DSType_Lambda,
		LambdaArn: props.FunctionArn,
	})
	return &impl{name: alias, resource: ds}
}
