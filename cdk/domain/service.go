package domain

import (
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/module/api"
	"github.com/cevixe/cdk/module/handler"
	"github.com/cevixe/cdk/service/sns"
)

type ServiceProps struct {
	App      string                 `field:"required"`
	Name     string                 `field:"required"`
	Api      api.ApiConfigProps     `field:"optional"`
	Handlers []handler.HandlerProps `field:"optional"`
}

func NewService(scope constructs.Construct, props *ServiceProps) module.Module {

	log.Printf("🥋 Cevixe Domain Service: %s\n", props.Name)
	mod := module.New(scope, module.Service, props.App, props.Name)

	api.ConfigApi(mod, &props.Api)

	advancedBusArn := mod.Import("core", export.AdvancedBusArn)
	advancedBus := sns.LoadTopic(mod, "advancedbus", advancedBusArn)

	standardBusArn := mod.Import("core", export.StandardBusArn)
	standardBus := sns.LoadTopic(mod, "standardbus", standardBusArn)

	for _, hdl := range props.Handlers {
		handler.NewHandler(mod, advancedBus, standardBus, &hdl)
	}

	return mod
}
