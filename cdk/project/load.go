package project

import (
	"log"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cevixe/cdk/app"
	"github.com/cevixe/cdk/domain"
	"github.com/cevixe/cdk/spec/v20221023"
)

func Load(scope constructs.Construct) {

	buffer := specBytes()
	base := &Base{}
	unmarshal(buffer, &base)

	switch base.Version {
	case "2022-10-23":

		file := &spec.File{}
		unmarshal(buffer, &file)
		project := file.Project

		switch project.Kind {
		case spec.Kind_App:
			app.Load20221023(scope, project.Name, &project.Properties)
		case spec.Kind_Domain:
			domain.Load20221023(scope, project.Name, &project.Properties)
		default:
			log.Fatalf("unsupported project kind: %v", project.Kind)
		}
	default:
		log.Fatalf("unsupported file version: %v", base.Version)
	}
}
