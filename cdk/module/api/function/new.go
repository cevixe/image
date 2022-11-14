package function

import (
	"fmt"
	"log"

	"github.com/cevixe/cdk/common/file"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func New(mod module.Module, alias string, props *FunctionProps) Function {

	requestLocation := fmt.Sprintf("%s/assets/function/%s/request.vtl", mod.Location(), alias)
	responseLocation := fmt.Sprintf("%s/assets/function/%s/response.vtl", mod.Location(), alias)

	if !file.Exists(requestLocation) {
		log.Fatalf("cannot locate request template for function: %s", requestLocation)
	}

	if !file.Exists(responseLocation) {
		log.Fatalf("cannot locate response template for function: %s", responseLocation)
	}

	requestTemplate := file.GetFileContent(requestLocation)
	responseTemplate := file.GetFileContent(responseLocation)

	fn := appsync.NewFunction(mod, alias, &appsync.FunctionProps{
		ApiId:            props.ApiId,
		DataSourceName:   props.DatasourceName,
		RequestTemplate:  requestTemplate,
		ResponseTemplate: responseTemplate,
	})

	return &impl{name: alias, resource: fn}
}
