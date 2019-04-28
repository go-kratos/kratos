package service

import (
	"bytes"
	"context"
	"time"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/dao"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Upload upload.
func (s *Service) Upload(c context.Context, fileName, fileType string, t time.Time, body []byte, bfs *conf.Bfs) (location string, err error) {
	if len(body) > bfs.MaxFileSize {
		err = ecode.FileTooLarge
		return
	}
	if location, err = dao.Upload(c, fileName, fileType, t.Unix(), bytes.NewReader(body), bfs); err != nil {
		log.Error("Upload error(%v)", err)
		return
	}
	return
}
