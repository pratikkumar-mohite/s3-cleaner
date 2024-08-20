package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *S3Client) getS3Bucket(bucket_name string) string {
	output, err := c.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic("unable to list buckets, " + err.Error())
	}
	for _, bucket := range output.Buckets {
		if aws.ToString(bucket.Name) == bucket_name {
			return bucket_name
		}
	}
	if bucket_name == "" {
		panic("bucket name is empty")
	}
	return ""
}

func (c *S3Client) checkVersioningStatus(bucket string) string {
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}

	result, err := c.Client.GetBucketVersioning(context.TODO(), input)
	if err != nil {
		panic("failed to get bucket versioning, " + err.Error())
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
			fmt.Printf("failed to get page, %v", err)
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

	return objects
}

func (c *S3Client) GetS3BucketObjects() []S3BucketObject {
	bucket := c.getS3Bucket(c.Bucket)
	output, err := c.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		panic("unable to list objects, " + err.Error())
	}
	objects := make([]S3BucketObject, len(output.Contents))
	if c.checkVersioningStatus(bucket) == "Enabled" {
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
		panic("unable to delete object version, " + err.Error())
	}
}

func (c *S3Client) UploadS3BucketObjects(object_file_path string) {
	file, err := os.Open(object_file_path)
	if err != nil {
		panic("Error opening file:" + err.Error())
	}
	defer file.Close()

	key := strings.Split(object_file_path, "/")[2]

	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		panic("Error uploading file:" + err.Error())
	}

	fmt.Printf("File %s uploaded successfully\n", key)
}
