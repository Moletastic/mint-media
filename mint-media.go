package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type MediaClientConfig struct {
	AccessKey     string
	SecretKey     string
	Region        string
	BucketName    string
	PresignExpire time.Duration
}

type MediaClient struct {
	SVC    *s3.S3
	Config *MediaClientConfig
}

func (c *MediaClient) GetObjectURI(key string) (string, error) {
	req, _ := c.SVC.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.Config.BucketName),
		Key:    aws.String(key),
	})
	url, err := req.Presign(c.Config.PresignExpire)
	if err != nil {
		return "", err
	}
	return url, nil
}

func New(c *MediaClientConfig) (*MediaClient, error) {
	creds := credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}
	cfg := aws.NewConfig().WithRegion(c.Region).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)
	client := &MediaClient{
		SVC:    svc,
		Config: c,
	}
	return client, nil
}
