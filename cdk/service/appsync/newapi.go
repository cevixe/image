package appsync

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsappsync"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"
)

type ApiProps struct {
	OIDCIssuer string `field:"required"`
}

func NewApi(mod module.Module, alias string, props *ApiProps) awsappsync.CfnGraphQLApi {

	name := naming.NewName(mod, naming.ResType_GraphQLApi, alias)

	cfnprops := &awsappsync.CfnGraphQLApiProps{
		Name:               name.Logical(),
		AuthenticationType: jsii.String("OPENID_CONNECT"),
		XrayEnabled:        jsii.Bool(true),
		OpenIdConnectConfig: &awsappsync.CfnGraphQLApi_OpenIDConnectConfigProperty{
			Issuer: jsii.String(props.OIDCIssuer),
		},
	}

	api := awsappsync.NewCfnGraphQLApi(mod.Resource(), name.Logical(), cfnprops)

	return api
}
