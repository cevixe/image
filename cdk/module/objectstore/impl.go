package objectstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/cevixe/cdk/module"
)

type impl struct {
	module   module.Module
	name     string
	resource awss3.Bucket
	presign  awslambda.Function
}

func (o *impl) Module() module.Module {
	return o.module
}

func (o *impl) Name() string {
	return o.name
}

func (o *impl) Resource() awss3.Bucket {
	return o.resource
}

func (o *impl) Presign() awslambda.Function {
	return o.presign
}
