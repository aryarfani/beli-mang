package aws

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	AccessKeyId     string
	SecretAccessKey string
	BucketName      string
	Region          string
}

var storage *s3.Client

func (conf *S3Config) Initialize() error {
	creds := credentials.NewStaticCredentialsProvider(conf.AccessKeyId, conf.SecretAccessKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return err
	}
	storage = s3.NewFromConfig(cfg)
	return nil
}

func (conf *S3Config) Upload(fileName string, file *multipart.FileHeader) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(storage)
	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	output, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(conf.BucketName),
		Key:    aws.String(fileName),
		Body:   f,
	})

	return output, err
}
