package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3Connection(cfg aws.Config) *s3.Client {
	client := s3.NewFromConfig(cfg)
	return client
}

func ListS3Buckets(c *s3.Client) {
	output, err := c.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic("unable to list buckets, " + err.Error())
	}
	for _, object := range output.Buckets {
		fmt.Printf("key=%s timestamp=%v\n", aws.ToString(object.Name), object.CreationDate)
	}
}