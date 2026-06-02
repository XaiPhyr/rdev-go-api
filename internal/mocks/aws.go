package mocks

import "context"

type AWSTest struct{}

func (a AWSTest) UploadToS3(ctx context.Context, key string, data []byte) (string, error) {
	return "https://mock-s3-url.com/" + key, nil
}

func NewTestAWSService() AWSTest {
	return AWSTest{}
}
