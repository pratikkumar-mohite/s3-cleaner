package cmd

import (
	"fmt"

	"github.com/pratikkumar-mohite/s3-cleanup/pkg/aws"
)

func setup() aws.S3Client {
	config := aws.AWSConnection(getFromEnv("AWS_PROFILE"))
	return aws.S3Connection(config)
}

func S3Cleanup() {
	s3Client := setup()
	bucket := getFromEnv("AWS_DELETE_S3_BUCKET")
	objects := s3Client.GetS3BucketObjects(bucket)
	for _, object := range objects {
		if object.ObjectName != "" {
			if object.ObjectDeleteMarker != "" {
				s3Client.DeleteS3BucketObjectVersion(bucket, object.ObjectName, object.ObjectDeleteMarker)
			}

			for _, version := range object.ObjectVersion {
				s3Client.DeleteS3BucketObjectVersion(bucket, object.ObjectName, version)
			}
			fmt.Printf("Delete Object: %v\n", object.ObjectName)
		}
	}
}
