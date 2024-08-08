package cmd

import (
	"fmt"

	"github.com/pratikkumar-mohite/s3-cleanup/pkg/aws"
)

func S3Cleanup() {
	config := aws.AWSConnection(getFromEnv("AWS_PROFILE"))
	s3Client := aws.S3Connection(config)
	objects := s3Client.GetS3BucketObjects(getFromEnv("AWS_DELETE_S3_BUCKET"))
	for _, object := range objects {
		if object.Object != "" {
			fmt.Println("Object Found: ", object.Object)
		}
	}
}