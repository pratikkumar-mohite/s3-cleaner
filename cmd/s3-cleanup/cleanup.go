package cmd

import (
	"fmt"
	"time"

	"github.com/pratikkumar-mohite/s3-cleanup/pkg/aws"
)

func setup() aws.S3Client {
	config := aws.AWSConnection(getFromEnv("AWS_PROFILE"))
	client := aws.S3Connection(config)
	client.Bucket = getFromEnv("AWS_DELETE_S3_BUCKET")
	return client
}

func S3Cleanup() {
	s3Client := setup()
	object1 := s3Client.UploadS3BucketObjects("test/files/file1.txt")
	s3Client.UploadS3BucketObjects("test/files/file2.txt")
	s3Client.UploadS3BucketObjects("test/files/file1.txt")
	s3Client.UploadS3BucketObjects("test/files/file2.txt")
	s3Client.DeleteS3BucketObjectVersion("file1.txt", object1)

	startTime := time.Now()

	objects := s3Client.GetS3BucketObjects()
	for _, object := range objects {
		if object.ObjectName != "" {
			if object.ObjectDeleteMarker != "" {
				s3Client.DeleteS3BucketObjectVersion(object.ObjectName, object.ObjectDeleteMarker)
			}

			for _, version := range object.ObjectVersion {
				s3Client.DeleteS3BucketObjectVersion(object.ObjectName, version)
			}
			fmt.Printf("Deleted Object: %v\n", object.ObjectName)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println("Total time taken for object deletion: ", elapsedTime)
}
