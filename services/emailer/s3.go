package emailer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader interface {
	Upload([]byte, string) (string, error) // string is the url of the uploaded content
}

type S3Uploader struct {
	u       *manager.Uploader
	viewURL string
}

func NewS3Config(endpoint string) (aws.Config, error) {
	var options []func(*config.LoadOptions) error

	// this endpoint is only used in integration tests and local dev
	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               endpoint,
					HostnameImmutable: true,
					SigningRegion:     "us-east-1",
				}, nil
			}

			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		options = append(options, config.WithEndpointResolverWithOptions(customResolver))
		options = append(options, config.WithRegion("us-east-1"))
		options = append(options, config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "minio",
				SecretAccessKey: "miniominio",
			}, nil
		})))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), options...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("load config: %w", err)
	}

	return cfg, nil
}

func NewS3Uploader(hostURL, viewURL string) *S3Uploader {
	cfg, err := NewS3Config(hostURL)
	if err != nil {
		log.Panic(err)
	}

	client := s3.NewFromConfig(cfg)
	up := manager.NewUploader(client)
	return &S3Uploader{
		u:       up,
		viewURL: viewURL,
	}
}

func (s S3Uploader) Upload(content []byte, name string) (string, error) {
	res, err := s.u.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String("mailer-emails"),
		Key:         aws.String(name),
		Body:        bytes.NewBuffer(content),
		ContentType: aws.String("text/html;charset=utf-8"),
	})
	if err != nil {
		return "", err
	}

	u, err := url.Parse(res.Location)
	if err != nil {
		return "", fmt.Errorf("could not parse url: %w", err)
	}

	u.Host = s.viewURL

	return u.String(), nil
}
