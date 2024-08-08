package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3Connection(cfg aws.Config) *s3.Client {
	client := s3.NewFromConfig(cfg)
	return client
}

func getS3Buckets(c *s3.Client) []S3Bucket {
	output, err := c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
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

func getS3Bucket(c *s3.Client, bucket_name string) string {
	buckets := getS3Buckets(c)
	for _, bucket := range buckets {
		if bucket.Name == bucket_name {
			return bucket.Name
		}
	}
	if bucket_name == "" {
		panic("bucket name is empty")
	}
	return ""
}

func GetS3BucketObjects(c *s3.Client, bucket_name string) []S3BucketObject {
	bucket := getS3Bucket(c, bucket_name)
	output, err := c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		panic("unable to list objects, " + err.Error())
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