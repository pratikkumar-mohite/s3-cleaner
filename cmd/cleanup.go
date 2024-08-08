package cmd

import (
	"github.com/pratikkumar-mohite/s3-cleanup/pkg/aws"
)

func S3Cleanup() {
	config := aws.AWSConnection(getFromEnv("AWS_PROFILE"))
	s3Client := aws.S3Connection(config)
	aws.ListS3Buckets(s3Client)
}