package service

import (
	"fmt"
	"io"

	"github.com/yangb8/webservice/common/config"
)

// Storage ...
type Storage interface {
	Upload(key string, body io.ReadCloser) error
	Download(key string) (io.ReadCloser, error)
	Delete(key string) error
}

func GetStorage(c *config.Config) (s Storage) {
	switch c.Storage.Kind {
	case "local":
		s = NewLocalStorage(c.Storage.Local.Path)
	case "s3":
		s3cfg := &c.Storage.S3
		s = NewS3Storage(s3cfg.URL, s3cfg.ID, s3cfg.Key, s3cfg.Bucket)
	default:
		panic(fmt.Sprintf("unsupported storage kind: %s", c.Storage.Kind))
	}
	return s
}
