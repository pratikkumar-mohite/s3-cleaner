package cleaner

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
		objectChan <- object
	}

	close(objectChan)
	wg.Wait()
}

func S3Cleanup(profile, region, bucket, prefix *string, deleteBucket, listObjects bool) {
	var s3Client aws.S3Client
	if *profile != "" && *region != "" && *bucket != "" {
		s3Client = setup(*profile, *region, *bucket)
		if *prefix != "" {
			s3Client.Prefix = *prefix
		}
	} else {
		s3Client = setup(getFromEnv("AWS_PROFILE"), getFromEnv("AWS_REGION"), getFromEnv("AWS_S3_BUCKET"))
		s3Client.Prefix = getFromEnv("AWS_S3_PREFIX")
	}

	if getFromEnv("AWS_S3_LIST_OBJECTS") == "true" && getFromEnv("AWS_S3_DELETE_BUCKET") == "true" {
		log.Fatalf("Cannot list objects and delete bucket at the same time, please unset one of the environment variable AWS_S3_LIST_OBJECTS or AWS_S3_DELETE_BUCKET")
	}

	if getFromEnv("AWS_UPLOAD_TEST_FILES") == "true" {
		s3Upload(s3Client)
	}

	startTime := time.Now()

	objects := s3Client.GetS3BucketObjects()

	if len(objects) == 0 {
		log.Infof("No objects found in bucket %s\n", s3Client.Bucket)
	} else if listObjects || getFromEnv("AWS_S3_LIST_OBJECTS") == "true" {
		for _, object := range objects {
			log.Infof("Object Name: %s, Object Version: %v, Object Delete Marker: %s\n", object.ObjectName, object.ObjectVersion, object.ObjectDeleteMarker)
		}
		log.Infof("%d objects found in bucket %s\n", len(objects), s3Client.Bucket)
	} else {
		concurrentCleanup(s3Client, objects)
	}

	elapsedTime := time.Since(startTime)

	if deleteBucket || getFromEnv("AWS_S3_DELETE_BUCKET") == "true" {
		s3Client.S3BucketDelete()
	}

	log.Infof("Time taken for bucket/objects deletion: %v", elapsedTime)
}
