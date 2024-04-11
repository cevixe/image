package s3

import (
	//"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/naming"
)

func NewBucket(mod module.Module, alias string, domain string) awss3.Bucket {

	name := naming.NewName(mod, naming.ResType_S3Bucket, alias)

	return awss3.NewBucket(
		mod.Resource(),
		name.Logical(),
		&awss3.BucketProps{
			BucketName:         jsii.String(domain),
			EnforceSSL:         jsii.Bool(true),
			Versioned:          jsii.Bool(false),
			PublicReadAccess:   jsii.Bool(false),
			EventBridgeEnabled: jsii.Bool(false),
			//AutoDeleteObjects:  jsii.Bool(true),
			//RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
			Cors: &[]*awss3.CorsRule{
				{
					AllowedMethods: &[]awss3.HttpMethods{
						awss3.HttpMethods_GET,
						awss3.HttpMethods_PUT,
					},
					AllowedOrigins: &[]*string{
						jsii.String("*"),
					},
					AllowedHeaders: &[]*string{
						jsii.String("*"),
					},
				},
			},
		},
	)
}
