package aws

import (
	"context"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *S3Client) getS3Bucket(bucket_name string) string {
	output, err := c.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Unable to list buckets, " + err.Error())
	}
	for _, bucket := range output.Buckets {
		if aws.ToString(bucket.Name) == bucket_name {
			return bucket_name
		}
	}
	if bucket_name == "" {
		log.Fatalf("Bucket name is empty %s", bucket_name)
	}
	return ""
}

func (c *S3Client) checkVersioningStatus(bucket string) string {
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}

	result, err := c.Client.GetBucketVersioning(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to get bucket versioning, %v" + err.Error())
	}

	return string(result.Status)
}

func (c *S3Client) listObjectVersions(bucket *string) []S3BucketObject {
	var objects []S3BucketObject
	var objectsMap = make(map[string]*S3BucketObject)
	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(*bucket),
	}

	paginator := s3.NewListObjectVersionsPaginator(c.Client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("Failed to get page, %v", err)
		}

		for _, version := range page.Versions {
			if _, exists := objectsMap[*version.Key]; !exists {
				objectsMap[*version.Key] = &S3BucketObject{ObjectName: *version.Key}
			}

			objectsMap[*version.Key].ObjectVersion = append(objectsMap[*version.Key].ObjectVersion, *version.VersionId)
		}

		for _, deleteMarker := range page.DeleteMarkers {
			if _, exists := objectsMap[*deleteMarker.Key]; !exists {
				objectsMap[*deleteMarker.Key] = &S3BucketObject{ObjectName: *deleteMarker.Key}
			}

			objectsMap[*deleteMarker.Key].ObjectDeleteMarker = *deleteMarker.VersionId
		}
	}
	for _, object := range objectsMap {
		objects = append(objects, *object)
	}

	if len(objects) == 0 {
		log.Infof("No objects found in bucket %s\n", *bucket)
		return objects
	}

	return objects
}

func (c *S3Client) GetS3BucketObjects() []S3BucketObject {
	bucket := c.getS3Bucket(c.Bucket)
	output, err := c.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		log.Fatalf("Unable to list objects, %v" + err.Error())
	}
	objects := make([]S3BucketObject, len(output.Contents))
	if c.checkVersioningStatus(bucket) == "Enabled" {
		log.Infof("Versioning is enabled for bucket %s\n", bucket)
		return c.listObjectVersions(&bucket)
	}
	for index, object := range output.Contents {
		key := aws.ToString(object.Key)
		if !strings.HasSuffix(key, "/") {
			objects[index] = S3BucketObject{
				ObjectName: key,
			}
		}
	}
	return objects
}

func (c *S3Client) DeleteS3BucketObjectVersion(object_name string, version_id string) {
	_, err := c.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket:    &c.Bucket,
		Key:       &object_name,
		VersionId: &version_id,
	})
	if err != nil {
		log.Errorf("Unable to delete object version, %v" + err.Error())
	}
}

func (c *S3Client) UploadS3BucketObjects(object_file_path string) string {
	file, err := os.Open(object_file_path)
	if err != nil {
		log.Fatalf("Error opening file: %v" + err.Error())
	}
	defer file.Close()

	key := strings.Split(object_file_path, "/")[2]

	object, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		log.Fatalf("Error uploading file: %v" + err.Error())
	}
	log.Infof("File %s uploaded successfully\n", key)
	return *object.VersionId
}

func (c *S3Client) S3BucketDelete() {
	_, err := c.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: &c.Bucket,
	})
	if err != nil {
		log.Fatalf("Unable to delete bucket, %v" + err.Error())
	}
	log.Infof("Bucket %s deleted successfully\n", c.Bucket)
}