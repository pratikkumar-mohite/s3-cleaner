package aws

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3Client struct {
	Client *s3.Client
	Bucket string
	Prefix string
}

type S3BucketObject struct {
	ObjectName         string
	ObjectVersion      []string
	ObjectDeleteMarker string
}
