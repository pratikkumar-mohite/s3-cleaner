package aws

import (
	"fmt"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3Connection(cfg aws.Config) S3Client {
	client := S3Client{
		Client: s3.NewFromConfig(cfg),
	}
	return client
}

func (c *S3Client)getS3Buckets() []S3Bucket {
	output, err := c.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic("unable to list buckets, " + err.Error())
	}
	buckets := make([]S3Bucket, len(output.Buckets))
	for index, object := range output.Buckets {
		buckets[index] = S3Bucket{
			Name: aws.ToString(object.Name),
		}
	}
	return buckets
}

func (c *S3Client)getS3Bucket(bucket_name string) *S3BucketOptions {
	buckets := c.getS3Buckets()
	for _, bucket := range buckets {
		if bucket.Name == bucket_name {
			return &S3BucketOptions{Name: bucket.Name, Versioning:  "Enabled"}
		}
	}
	if bucket_name == "" {
		panic("bucket name is empty")
	}
	return &S3BucketOptions{Name: "", Versioning:  ""}
}

func (c *S3Client)checkVersioning(bucket string) string {
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}

	result, err := c.Client.GetBucketVersioning(context.TODO(), input)
	if err != nil {
		panic("failed to get bucket versioning, " + err.Error())
	}

	return string(result.Status)
}

func (c *S3Client)listObjectVersions(bucket *string) {
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
			if *version.Key != "" {
				fmt.Printf("Object: %s, Version ID: %s\n", *version.Key, *version.VersionId)
			}
		}
		for _, marker := range page.DeleteMarkers {
			fmt.Printf("Delete Marker Object: %s, Version ID: %s\n", *marker.Key, *marker.VersionId)
		}
	}
}

func (c *S3Client)GetS3BucketObjects(bucket_name string) []S3BucketObject {
	bucket := c.getS3Bucket(bucket_name)
	output, err := c.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket.Name,
	})
	if err != nil {
		panic("unable to list objects, " + err.Error())
	}
	if c.checkVersioning(bucket.Name) == "Enabled" {
		c.listObjectVersions(&bucket.Name)
	}
	objects := make([]S3BucketObject, len(output.Contents))
	for index, object := range output.Contents {
		key := aws.ToString(object.Key)
		if !strings.HasSuffix(key, "/") {
			objects[index] = S3BucketObject{
				Object: key,
			}
		}
	}
	return objects
}