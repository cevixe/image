package domain

import (
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/api"
	"github.com/cevixe/cdk/module/api/datasource"
	"github.com/cevixe/cdk/module/api/function"
	"github.com/cevixe/cdk/module/objectstore"
	"github.com/cevixe/cdk/module/statestore"
	"github.com/cevixe/cdk/service/iam"
	"github.com/cevixe/cdk/service/route53"
)

type PlatformProps struct {
	App     string   `field:"required"`
	Name    string   `field:"required"`
	Indexes []string `field:"optional"`
}

func NewPlatform(scope constructs.Construct, props *PlatformProps) module.Module {

	log.Printf("🥋 Cevixe Domain Platform: %s\n", props.Name)
	mod := module.New(scope, module.Platform, props.App, props.Name)

	advancedBusArn := mod.Import("core", export.AdvancedBusArn)
	standardBusArn := mod.Import("core", export.StandardBusArn)

	statestore := statestore.NewStateStore(mod, "statestore",
		&statestore.StateStoreProps{
			Indexes:        props.Indexes,
			AdvancedBusArn: advancedBusArn,
			StandardBusArn: standardBusArn,
		},
	)
	mod.Export(export.StateStoreName, *statestore.Resource().TableName())
	mod.Export(export.StateStoreArn, *statestore.Resource().TableArn())

	zone := route53.LoadZone(mod, "zone", &route53.ZoneProps{
		ID:   mod.Import("core", export.HostedZoneId),
		Name: mod.Import("core", export.HostedZoneName),
	})

	objectstore := objectstore.NewObjectStore(mod,
		&objectstore.ObjectStoreProps{Alias: "objectstore", Zone: zone})

	gqlapi := api.NewApi(mod, props.Name,
		&api.ApiProps{
			Zone:        zone,
			ObjectStore: objectstore,
			StateStore:  statestore,
		},
	)
	mod.Export(export.GraphQLApiId, *gqlapi.Resource().AttrApiId())
	mod.Export(export.GraphQLApiArn, *gqlapi.Resource().AttrArn())
	mod.Export(export.GraphQLApiUrl, gqlapi.URL())
	mod.Export(export.GraphQLApiKey, *gqlapi.Key().AttrApiKey())

	gqlapi.Role().AddToPrincipalPolicy(
		iam.NewDynCrudPol(*statestore.Resource().TableArn()))
	gqlapi.Role().AddToPrincipalPolicy(
		iam.NewLambdaInvokePol(*objectstore.Presign().FunctionArn()))

	configureDefaultFunctions(mod, gqlapi, statestore, objectstore)

	return mod
}

func configureDefaultFunctions(
	mod module.Module,
	api api.Api,
	statestore statestore.StateStore,
	objectstore objectstore.ObjectStore,
) {

	datasource.NewDataSource(mod, "mock", &datasource.DataSourceProps{
		ApiId: *api.Resource().AttrApiId(),
		Type:  datasource.DSType_Mock,
	})

	statestoreds := datasource.NewDataSource(mod, "statestore", &datasource.DataSourceProps{
		ApiId:      *api.Resource().AttrApiId(),
		Type:       datasource.DSType_Table,
		RoleArn:    *api.Role().RoleArn(),
		Region:     *statestore.Resource().Stack().Region(),
		StateStore: statestore.Resource(),
	})

	objectstoreds := datasource.NewDataSource(mod, "objectstore", &datasource.DataSourceProps{
		ApiId:       *api.Resource().AttrApiId(),
		Type:        datasource.DSType_Lambda,
		RoleArn:     *api.Role().RoleArn(),
		FunctionArn: *objectstore.Presign().FunctionArn(),
	})

	function.NewStateStoreCreateFn(mod, &function.FunctionProps{
		ApiId:          *statestoreds.Resource().ApiId(),
		DatasourceName: statestoreds.Name(),
	}).Resource().AddDependsOn(statestoreds.Resource())
	function.NewStateStoreUpdateFn(mod, &function.FunctionProps{
		ApiId:          *statestoreds.Resource().ApiId(),
		DatasourceName: statestoreds.Name(),
	}).Resource().AddDependsOn(statestoreds.Resource())
	function.NewStateStoreDeleteFn(mod, &function.FunctionProps{
		ApiId:          *statestoreds.Resource().ApiId(),
		DatasourceName: statestoreds.Name(),
	}).Resource().AddDependsOn(statestoreds.Resource())
	function.NewStateStoreFindOneFn(mod, &function.FunctionProps{
		ApiId:          *statestoreds.Resource().ApiId(),
		DatasourceName: statestoreds.Name(),
	}).Resource().AddDependsOn(statestoreds.Resource())
	function.NewStateStoreFindAllFn(mod, &function.FunctionProps{
		ApiId:          *statestoreds.Resource().ApiId(),
		DatasourceName: statestoreds.Name(),
	}).Resource().AddDependsOn(statestoreds.Resource())
	function.NewStateStoreFindByFn(mod, &function.FunctionProps{
		ApiId:          *statestoreds.Resource().ApiId(),
		DatasourceName: statestoreds.Name(),
	}).Resource().AddDependsOn(statestoreds.Resource())

	function.NewObjectStoreUploadFn(mod, &function.FunctionProps{
		ApiId:          *objectstoreds.Resource().ApiId(),
		DatasourceName: objectstoreds.Name(),
	}).Resource().AddDependsOn(objectstoreds.Resource())
	function.NewObjectStoreDownloadFn(mod, &function.FunctionProps{
		ApiId:          *objectstoreds.Resource().ApiId(),
		DatasourceName: objectstoreds.Name(),
	}).Resource().AddDependsOn(objectstoreds.Resource())
}
