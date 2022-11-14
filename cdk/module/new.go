package module

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ModuleProps struct {
	App  string `field:"optional"`
	Type Type   `field:"required"`
}

func New(scope constructs.Construct, typ Type, app string, name string) Module {

	var location string
	switch typ {
	case Service:
		location = os.Getenv("CEVIXE_MOD_HOME")
	case Platform:
		location = os.Getenv("CEVIXE_APP_HOME")
	default:
		log.Fatalf("unsupported module type: %v", typ)
	}

	stackName := fmt.Sprintf("cvx-%s-%s-%s", app, name, string(typ))
	resource := awscdk.NewStack(scope, jsii.String(name), &awscdk.StackProps{
		StackName: jsii.String(stackName),
		Env: &awscdk.Environment{
			Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
			Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
		},
	})

	return &impl{app: app, name: name, location: location, resource: resource}
}
