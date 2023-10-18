package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/common/file"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/api/apollo"
	"github.com/cevixe/cdk/module/api/datasource"
	"github.com/cevixe/cdk/module/api/function"
	"github.com/cevixe/cdk/module/api/resolver"
	lambda "github.com/cevixe/cdk/module/function"
	"github.com/cevixe/cdk/service/appsync"
	"github.com/cevixe/cdk/service/iam"
)

func configLambdaDS(mod module.Module, fn lambda.Function) {

	commandStoreArn := mod.Import("core", export.CommandStoreArn)
	commandStoreName := mod.Import("core", export.CommandStoreName)

	stateStoreArn := mod.Import(mod.Name(), export.StateStoreArn)
	stateStoreName := mod.Import(mod.Name(), export.StateStoreName)

	objectStoreArn := mod.Import(mod.Name(), export.ObjectStoreArn)
	objectStoreName := mod.Import(mod.Name(), export.ObjectStoreName)

	fn.Resource().AddEnvironment(jsii.String("CVX_STATE_STORE"), jsii.String(stateStoreName), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_OBJECT_STORE"), jsii.String(objectStoreName), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_COMMAND_STORE"), jsii.String(commandStoreName), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_APP_NAME"), jsii.String(mod.App()), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_DOMAIN_NAME"), jsii.String(mod.Name()), nil)
	fn.Resource().AddEnvironment(jsii.String("CVX_HANDLER_NAME"), jsii.String(fn.Name()), nil)

	fn.Resource().AddToRolePolicy(iam.NewDynReadPol("*"))
	fn.Resource().AddToRolePolicy(iam.NewDynCrudPol(stateStoreArn))
	fn.Resource().AddToRolePolicy(iam.NewS3CrudPol(objectStoreArn))
	fn.Resource().AddToRolePolicy(iam.NewDynWritePol(commandStoreArn))
}

func ConfigApi(mod module.Module, props *ApiConfigProps) {

	schemaLocation := findSchemaLocation(mod)
	schemaContent := file.GetFileContent(schemaLocation)
	apiId := mod.Import(mod.Name(), export.GraphQLApiId)
	stateStoreName := mod.Import(mod.Name(), export.StateStoreName)

	schema := appsync.NewSchema(mod, mod.Name(), &appsync.SchemaProps{
		ApiId:      apiId,
		Definition: schemaContent,
	})

	apollo.NewEntitiesResolver(mod, &apollo.EntitiesResolverProps{
		ApiId:          apiId,
		StateStoreName: stateStoreName,
	})
	apollo.NewServiceResolver(mod, &apollo.ServiceResolverProps{
		ApiId:  apiId,
		Schema: schemaContent,
	})

	functionsMap := map[string]string{
		"sscreatefn":   mod.Import(mod.Name(), export.StateStoreCreateFn),
		"ssupdatefn":   mod.Import(mod.Name(), export.StateStoreUpdateFn),
		"ssdeletefn":   mod.Import(mod.Name(), export.StateStoreDeleteFn),
		"ssfindallfn":  mod.Import(mod.Name(), export.StateStoreFindAllFn),
		"ssfindonefn":  mod.Import(mod.Name(), export.StateStoreFindOneFn),
		"ssfindbyfn":   mod.Import(mod.Name(), export.StateStoreFindByFn),
		"osuploadfn":   mod.Import(mod.Name(), export.ObjectStoreUploadFn),
		"osdownloadfn": mod.Import(mod.Name(), export.ObjectStoreDownloadFn),
	}

	datasourceMap := map[string]datasource.DataSource {}

	for _, item := range props.DataSources {
		if item.Type != DataSourceType_Lambda { continue }
		entryFormat := "/cmd/datasource/%s"
		entry := fmt.Sprintf(entryFormat, item.Name)
		fn := lambda.NewFunction(mod, item.Name, entry)
		configLambdaDS(mod, fn)

		ds := datasource.NewDataSource(mod, item.Name, &datasource.DataSourceProps{
			ApiId:       apiId,
			Type:        datasource.DSType_Lambda,
			RoleArn:     mod.Import(mod.Name(), export.GraphQLApiRole),
			FunctionArn: *fn.Resource().FunctionArn(),
		})
		datasourceMap[item.Name] = ds
	}

	for _, item := range props.Functions {
		fn := function.New(mod, item.Name, &function.FunctionProps{
			ApiId:          apiId,
			DatasourceName: item.DataSource,
		})
		if datasourceMap[item.DataSource] != nil {
			fn.Resource().AddDependsOn(datasourceMap[item.DataSource].Resource())
		}
		
		functionsMap[item.Name] = *fn.Resource().AttrFunctionId()
	}

	for _, item := range props.Resolvers {
		operationItems := strings.Split(item.Operation, "/")
		functionIds := make([]string, 0)
		for _, name := range item.Functions {
			if functionsMap[name] == "" {
				log.Fatalf("api function not found: %s\n", name)
			}
			functionIds = append(functionIds, functionsMap[name])
		}
		resolver.NewResolver(mod, item.Name, &resolver.ResolverProps{
			ApiId:     apiId,
			Type:      operationItems[0],
			Field:     operationItems[1],
			Functions: functionIds,
			Schema:    schema,
		})
	}
}

func findSchemaLocation(mod module.Module) string {

	options := []string{
		fmt.Sprintf("%s/api/schema.gql", mod.Location()),
		fmt.Sprintf("%s/api/schema.graphql", mod.Location()),
	}

	for _, location := range options {
		if file.Exists(location) {
			return location
		}
	}

	log.Fatalf("cannot find graphql schema: %s", options)
	return ""
}
