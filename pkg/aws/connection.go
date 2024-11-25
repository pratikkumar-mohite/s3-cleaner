package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

func AWSConnection(profile, region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	cfg.Region = region
	return cfg
}

func S3Connection(cfg aws.Config) S3Client {
	client := S3Client{
		Client: s3.NewFromConfig(cfg),
	}
	return client
}
