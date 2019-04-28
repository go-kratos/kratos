package service

import (
	"context"
	"io"
)

// Upload http upload file.
func (s *Service) Upload(c context.Context, fileName, fileType string, expire int64, body io.Reader) (location string, err error) {
	return s.bfs.Upload(c, fileName, fileType, expire, body)
}
