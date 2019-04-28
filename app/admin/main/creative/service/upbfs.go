package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	"go-common/app/admin/main/creative/model/bfs"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Upload to bfs
func (s *Service) Upload(c context.Context, fileName string, fileType string, timing int64, body []byte) (location string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		return
	}
	if len(body) > s.conf.Bfs.MaxFileSize {
		err = ecode.FileTooLarge
		return
	}
	if location, err = s.dao.Upload(c, fileName, fileType, timing, body, s.conf.Bfs); err != nil {
		log.Error("s.upload.Upload() error(%v)", err)
	}
	return
}

// ParseFile analyses file info
func (s *Service) ParseFile(c context.Context, content []byte) (file *bfs.FileInfo, err error) {
	fType := http.DetectContentType(content)
	// file md5
	md5hash := md5.New()
	if _, err = io.Copy(md5hash, bytes.NewReader(content)); err != nil {
		log.Error("resource uploadFile.Copy error(%v)", err)
		return
	}
	md5 := md5hash.Sum(nil)
	fMd5 := hex.EncodeToString(md5[:])
	file = &bfs.FileInfo{
		Md5:  fMd5,
		Type: fType,
		Size: int64(len(content)),
	}
	return
}
