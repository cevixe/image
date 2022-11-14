module github.com/cevixe/app

go 1.19

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go-v2/config v1.17.9
	github.com/aws/aws-sdk-go-v2/service/sns v1.18.2
	github.com/cevixe/sdk v0.0.0
	github.com/stoewer/go-strcase v1.2.0 // indirect
)

replace github.com/cevixe/sdk => ../sdk

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.17.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.29.1
	github.com/aws/jsii-runtime-go v1.70.0
)

require (
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/aws/aws-sdk-go-v2 v1.17.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.22 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodbstreams/attributevalue v1.10.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.13.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.17.0 // indirect
	github.com/aws/smithy-go v1.13.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/tools v0.1.12 // indirect
)
