package aws

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSService interface {
	UploadToS3(ctx context.Context, key string, data []byte) (string, error)
}

type service struct {
	s3Client *s3.Client
	bucket   string
	region   string
}

func NewAWSService(region, bucket string) (*service, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &service{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
		region:   region,
	}, nil
}

func (s *service) UploadToS3(ctx context.Context, key string, data []byte) (string, error) {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(http.DetectContentType(data)),
	})

	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)

	return publicURL, nil
}
