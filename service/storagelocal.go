package service

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

type StorageLocal struct {
}

func NewLocalStorage(path string) Storage {
	return &StorageLocal{}
}

func (s *StorageLocal) Upload(key string, body io.ReadCloser) error {
	f, err := os.Create("/tmp/test/" + key)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err = w.ReadFrom(body); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func (s *StorageLocal) Download(key string) (io.ReadCloser, error) {
	f, err := os.Open("/tmp/test/" + key)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bufio.NewReader(f)), nil
}

func (s *StorageLocal) Delete(key string) error {
	if err := os.Remove("/tmp/test/" + key); err != nil {
		return err
	}
	return nil
}

var _ Storage = &StorageLocal{}
