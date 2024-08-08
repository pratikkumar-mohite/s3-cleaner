package aws

import (
	"context"

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

func GetS3Bucket(c *s3.Client, bucket_name string) string {
	buckets := getS3Buckets(c)
	for _, bucket := range buckets {
		if bucket.Name == bucket_name {
			return bucket.Name
		}
	}
	return ""
}