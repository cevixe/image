package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/project"
)

func main() {
	defer jsii.Close()
	scope := awscdk.NewApp(nil)
	project.Load(scope)
	scope.Synth(nil)
}
