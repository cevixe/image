package resolver

import (
	"fmt"
	"log"

	"github.com/cevixe/cdk/common/file"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/appsync"
)

func NewResolver(mod module.Module, alias string, props *ResolverProps) Resolver {

	requestLocation := fmt.Sprintf("%s/assets/resolver/%s/request.vtl", mod.Location(), alias)
	responseLocation := fmt.Sprintf("%s/assets/resolver/%s/response.vtl", mod.Location(), alias)

	if !file.Exists(requestLocation) {
		log.Fatalf("cannot locate request template for resolver: %s", requestLocation)
	}

	if !file.Exists(responseLocation) {
		log.Fatalf("cannot locate response template for resolver: %s", responseLocation)
	}

	requestTemplate := file.GetFileContent(requestLocation)
	responseTemplate := file.GetFileContent(responseLocation)

	resource := appsync.NewResolver(mod, alias, &appsync.ResolverProps{
		ApiId:            props.ApiId,
		Type:             props.Type,
		Field:            props.Field,
		Functions:        props.Functions,
		RequestTemplate:  requestTemplate,
		ResponseTemplate: responseTemplate,
	})

	resource.AddDependsOn(props.Schema)

	return &impl{
		name:     alias,
		resource: resource,
	}
}
