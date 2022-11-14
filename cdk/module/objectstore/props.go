package objectstore

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
)

type ObjectStoreProps struct {
	Alias string                 `field:"required"`
	Zone  awsroute53.IHostedZone `field:"required"`
}
