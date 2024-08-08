package cmd

import (
	"fmt"

	"github.com/pratikkumar-mohite/s3-cleanup/pkg/aws"
)

func S3Cleanup() {
	config := aws.AWSConnection(getFromEnv("AWS_PROFILE"))
	s3Client := aws.S3Connection(config)
	bucket := aws.GetS3Bucket(s3Client, getFromEnv("AWS_DELETE_S3_BUCKET"))
	if bucket == "" {
		fmt.Println("Bucket Not Found")
		return
	}
	fmt.Println("Bucket Found: ", bucket)
}