package app

import (
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/core"
	"github.com/cevixe/cdk/domain"
)

type ApplicationProps struct {
	Domains []ApplicationProps_Domain
}
type ApplicationProps_Domain struct {
	Name    string
	Indexes []string
}

func NewApplication(scope constructs.Construct, app string, props *ApplicationProps) {

	log.Printf("🥋 Cevixe App Platform: %s\n", app)
	corePlatform := core.NewPlatform(scope, app)
	for _, dom := range props.Domains {
		domainPlatform := domain.NewPlatform(scope, &domain.PlatformProps{
			App:     app,
			Name:    dom.Name,
			Indexes: dom.Indexes,
		})
		domainPlatform.DependsOn(corePlatform)
	}
}
