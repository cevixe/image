package naming

type ResType string

// AWS AppSync resource types
const (
	ResType_GraphQLDomain     ResType = "gdo"
	ResType_GraphQLDomainLink ResType = "gdl"
	ResType_GraphQLApi        ResType = "gap"
	ResType_GraphQLApiKey     ResType = "gak"
	ResType_GraphQLSchema     ResType = "gsc"
	ResType_GraphQLResolver   ResType = "grs"
	ResType_GraphQLFunction   ResType = "gfn"
	ResType_GraphQLDataSource ResType = "gds"
)

// AWS Certificate Manager
const (
	ResType_ACMCertificate ResType = "cer"
)

// AWS Route53 DNS Service
const (
	ResType_Route53Zone   ResType = "dns"
	ResType_Route53Record ResType = "rec"
)

// AWS Identity & Access Management
const (
	ResType_IAMRole   ResType = "rol"
	ResType_IAMPolicy ResType = "pol"
)

// AWS Lambda
const (
	ResType_Lambda ResType = "fun"
)

// AWS Simple Notification Service
const (
	ResType_SNSTopic        ResType = "top"
	ResType_SNSSubscription ResType = "sub"
)

// AWS Simple Queue Service
const (
	ResType_SQSQueue       ResType = "que"
	ResType_SQSQueuePolicy ResType = "qup"
)

// AWS DynamoDB
const (
	ResType_DynamodbTable ResType = "dyn"
)

// AWS Simple Storage Service (S3)
const (
	ResType_S3Bucket ResType = "s3b"
)
