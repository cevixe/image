package objectstore

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/app/pkg/location"
	"github.com/cevixe/cdk/common/export"
	"github.com/cevixe/cdk/module"
	"github.com/cevixe/cdk/service/iam"
	"github.com/cevixe/cdk/service/lambda"
	"github.com/cevixe/cdk/service/s3"
)

func NewObjectStore(mod module.Module, props *ObjectStoreProps) ObjectStore {

	bucketDomain := fmt.Sprintf("%s.%s.%s", props.Alias, mod.Name(), *props.Zone.ZoneName())
	bucket := s3.NewBucket(mod, props.Alias, bucketDomain)

	presign := lambda.NewGolangFunction(mod, "presign", location.ObjectPresign)
	presign.AddEnvironment(jsii.String("CVX_OBJECT_STORE"), bucket.BucketName(), nil)
	presign.AddToRolePolicy(iam.NewS3CrudPol(*bucket.BucketArn()))

	mod.Export(export.ObjectStoreName, *bucket.BucketName())
	mod.Export(export.ObjectStoreArn, *bucket.BucketArn())

	return &impl{
		module:   mod,
		name:     bucketDomain,
		resource: bucket,
		presign:  presign,
	}
}
