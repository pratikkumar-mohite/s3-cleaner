package main

import (
	"runtime"
	"sync"
	"time"

	"github.com/pratikkumar-mohite/s3-cleaner/pkg/aws"
	log "github.com/sirupsen/logrus"
)

func setup(profile, region, bucket string) aws.S3Client {
	config := aws.AWSConnection(profile, region)
	client := aws.S3Connection(config)
	client.Bucket = bucket
	return client
}

func s3Upload(s3Client aws.S3Client) {
	object1 := s3Client.UploadS3BucketObjects("test/files/file1.txt")
	s3Client.UploadS3BucketObjects("test/files/file2.txt")
	s3Client.UploadS3BucketObjects("test/files/file1.txt")
	s3Client.UploadS3BucketObjects("test/files/file2.txt")
	s3Client.DeleteS3BucketObjectVersion("file1.txt", object1)
}

func concurrentCleanup(s3Client aws.S3Client, objects []aws.S3BucketObject) {
	workerCount := runtime.NumCPU() * 2
	var wg sync.WaitGroup
	objectChan := make(chan aws.S3BucketObject, len(objects))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for object := range objectChan {
				if object.ObjectName != "" {
					if object.ObjectDeleteMarker != "" {
						s3Client.DeleteS3BucketObjectVersion(object.ObjectName, object.ObjectDeleteMarker)
					}

					var versionWG sync.WaitGroup

					if len(object.ObjectVersion) == 0 {
						versionWG.Add(1)
						go func(object_name string) {
							defer versionWG.Done()
							s3Client.DeleteS3BucketObject(object_name)
						}(object.ObjectName)
					}

					for _, version := range object.ObjectVersion {
						versionWG.Add(1)
						go func(version string) {
							defer versionWG.Done()
							s3Client.DeleteS3BucketObjectVersion(object.ObjectName, version)
						}(version)
					}
					versionWG.Wait()
				}
			}
		}()
	}
	for _, object := range objects {
		wg.Add(1)
		objectChan <- object
	}

	close(objectChan)
	wg.Wait()
}

func s3Cleanup(profile, region, bucket, prefix *string, bucketDelete bool) {
	var s3Client aws.S3Client
	if *profile != "" && *region != "" && *bucket != "" {
		s3Client = setup(*profile, *region, *bucket)
	} else {
		s3Client = setup(getFromEnv("AWS_PROFILE"), getFromEnv("AWS_REGION"), getFromEnv("AWS_S3_BUCKET"))
	}

	if *prefix != "" {
		s3Client.Prefix = *prefix
	} else {
		s3Client.Prefix = getFromEnv("AWS_S3_PREFIX")
	}

	if getFromEnv("AWS_UPLOAD_TEST_FILES") == "true" {
		s3Upload(s3Client)
	}

	startTime := time.Now()

	objects := s3Client.GetS3BucketObjects()

	if len(objects) == 0 {
		log.Infof("No objects found in bucket %s\n", s3Client.Bucket)
	} else {
		concurrentCleanup(s3Client, objects)
	}

	elapsedTime := time.Since(startTime)

	if bucketDelete || getFromEnv("AWS_S3_BUCKET_DELETE") == "true" {
		s3Client.S3BucketDelete()
	}

	log.Infof("Time taken for bucket/objects deletion: %v", elapsedTime)
}
