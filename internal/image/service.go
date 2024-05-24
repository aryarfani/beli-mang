package image

import (
	"beli-mang/pkg/aws"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	Upload(image *multipart.FileHeader) (imageUrl string, err error)
	ValidateFile(fileHeader *multipart.FileHeader) error
}

type service struct {
	s3Config *aws.S3Config
}

func NewService(s3Config *aws.S3Config) Service {
	return &service{
		s3Config: s3Config,
	}
}

func (s *service) Upload(image *multipart.FileHeader) (imageUrl string, err error) {
	newFileName := uuid.NewString() + ".jpg"
	output, err := s.s3Config.Upload(newFileName, image)
	if err != nil {
		return "", err
	}

	return output.Location, nil
}

func (s *service) ValidateFile(fileHeader *multipart.FileHeader) error {
	const maxFileSize = 2 * 1024 * 1024 // 2MB
	const minFileSize = 10 * 1024       // 10KB

	// Validate file size
	if fileHeader.Size > maxFileSize || fileHeader.Size < minFileSize {
		return errors.New("file size must be between 10KB and 2MB")
	}

	// Validate file extension
	ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
	if ext != "jpg" && ext != "jpeg" {
		return errors.New("only .jpg and .jpeg formats are allowed")
	}

	return nil
}
