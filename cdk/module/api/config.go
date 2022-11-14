package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/common/file"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/api/apollo"
	"github.com/cevixe/cdk/module/api/function"
	"github.com/cevixe/cdk/module/api/resolver"
	"github.com/cevixe/cdk/service/appsync"
)

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

	for _, item := range props.Functions {
		fn := function.New(mod, item.Name, &function.FunctionProps{
			ApiId:          apiId,
			DatasourceName: item.DataSource,
		})
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
