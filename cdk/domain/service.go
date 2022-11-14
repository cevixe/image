package domain

import (
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/api"
)

type ServiceProps struct {
	App  string             `field:"required"`
	Name string             `field:"required"`
	Api  api.ApiConfigProps `field:"optional"`
}

func NewService(scope constructs.Construct, props *ServiceProps) module.Module {

	log.Printf("🥋 Cevixe Domain Service: %s\n", props.Name)
	mod := module.New(scope, module.Service, props.App, props.Name)

	api.ConfigApi(mod, &props.Api)

	return mod
}
