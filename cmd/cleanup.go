package cmd

import (
	"fmt"

	"github.com/pratikkumar-mohite/s3-cleanup/pkg/aws"
)

func S3Cleanup() {
	config := aws.AWSConnection(getFromEnv("AWS_PROFILE"))
	s3Client := aws.S3Connection(config)
	objects := aws.GetS3BucketObjects(s3Client, getFromEnv("AWS_DELETE_S3_BUCKET"))
	fmt.Println("Objects Found: ", objects)
}