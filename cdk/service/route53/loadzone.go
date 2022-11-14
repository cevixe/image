package route53

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
)

type ZoneProps struct {
	ID   string
	Name string
}

func LoadZone(mod module.Module, alias string, props *ZoneProps) awsroute53.IHostedZone {

	return awsroute53.PublicHostedZone_FromHostedZoneAttributes(
		mod.Resource(), jsii.String(alias), &awsroute53.HostedZoneAttributes{
			HostedZoneId: jsii.String(props.ID),
			ZoneName:     jsii.String(props.Name),
		},
	)
}
