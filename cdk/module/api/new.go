package api

import (
	"fmt"

	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/acm"
	"github.com/cevixe/cdk/service/appsync"
	"github.com/cevixe/cdk/service/route53"
)

func NewApi(mod module.Module, alias string, props *ApiProps) Api {

	domainName := fmt.Sprintf("%s.%s", mod.Name(), *props.Zone.ZoneName())

	api := appsync.NewApi(mod, alias, &appsync.ApiProps{OIDCIssuer: props.OIDCIssuer})
	key := appsync.NewApiKey(mod, alias, &appsync.ApiKeyProps{Api: api})

	role := appsync.NewApiRole(mod, alias)

	certificate := acm.NewCertificate(mod, alias, &acm.CertificateProps{
		Zone:   props.Zone,
		Domain: domainName,
	})

	apiDomain := appsync.NewApiDomain(mod, alias, &appsync.ApiDomainProps{
		Api:         api,
		Domain:      domainName,
		Certificate: certificate,
	})

	route53.NewCnameRecord(mod, alias, &route53.CnameRecordProps{
		Zone:   props.Zone,
		Record: mod.Name(),
		Domain: *apiDomain.AttrAppSyncDomainName(),
	})

	return &apiImpl{
		module:   mod,
		name:     alias,
		record:   mod.Name(),
		domain:   *props.Zone.ZoneName(),
		key:      key,
		role:     role,
		resource: api,
	}
}
