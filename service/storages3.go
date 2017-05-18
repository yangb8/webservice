package service

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type StorageS3 struct {
	Client     *s3.S3
	Bucket     *string
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
}

func NewS3Storage(url, id, key, bucket string) Storage {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(id, key, ""),
		Endpoint:    aws.String(url),
		Region:      aws.String("us-west-2"), // TODO: make it configurable
	}
	sess := session.Must(session.NewSession(s3Config))

	c := s3.New(sess)

	// make sure Bucket is created
	if _, err := c.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
			case s3.ErrCodeBucketAlreadyOwnedByYou:
			default:
				panic(aerr)
			}
		} else {
			panic(err)
		}
	}

	return &StorageS3{
		Client:     c,
		Bucket:     aws.String(bucket),
		Downloader: s3manager.NewDownloader(sess),
		Uploader:   s3manager.NewUploader(sess),
	}
}

func (s *StorageS3) Upload(key string, body io.ReadCloser) error {
	_, err := s.Uploader.Upload(
		&s3manager.UploadInput{
			Bucket: s.Bucket,
			Key:    aws.String(key),
			Body:   body,
		})
	return err
}

func (s *StorageS3) Download(key string) (io.ReadCloser, error) {
	resp, err := s.Client.GetObject(
		&s3.GetObjectInput{
			Bucket: s.Bucket,
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (s *StorageS3) Delete(key string) error {
	_, err := s.Client.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: s.Bucket,
			Key:    aws.String(key),
		})
	return err
}

var _ Storage = &StorageS3{}
