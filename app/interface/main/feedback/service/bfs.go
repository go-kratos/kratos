package service

import (
	"bytes"
	"context"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

// Upload upload.
func (s *Service) Upload(c context.Context, fileName, fileType string, t time.Time, body []byte) (location string, err error) {
	if len(body) == 0 {
		err = ecode.FeedbackBodyNotExist
		return
	}
	if len(body) > s.c.Bfs.MaxFileSize {
		err = ecode.FeedbackBodyTooLarge
		return
	}
	if location, err = s.dao.Upload(c, fileName, fileType, t.Unix(), bytes.NewReader(body)); err != nil {
		log.Error("s.dao.Upload error(%v)", err)
		return
	}
	return
}
