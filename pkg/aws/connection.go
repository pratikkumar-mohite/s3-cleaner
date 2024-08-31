package aws

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func AWSConnection(profile string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	return cfg
}

func S3Connection(cfg aws.Config) S3Client {
	client := S3Client{
		Client: s3.NewFromConfig(cfg),
	}
	return client
}
