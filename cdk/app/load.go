package app

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/spec/v20221023"
)

func Load20221023(scope constructs.Construct, name string, props *spec.Properties) {

	domains := make([]ApplicationProps_Domain, 0)
	for _, dom := range props.Domains {
		domains = append(domains, ApplicationProps_Domain{
			Name:    dom.Name,
			Indexes: dom.Indexes,
		})
	}
	NewApplication(scope, name, &ApplicationProps{Domains: domains})
}
