package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	"go-common/app/job/main/appstatic/model"
	"go-common/library/log"
)

// ParseFile analyses file info
func (s *Service) ParseFile(content []byte) (file *model.FileInfo, err error) {
	fType := http.DetectContentType(content)
	// file md5
	md5hash := md5.New()
	if _, err = io.Copy(md5hash, bytes.NewReader(content)); err != nil {
		log.Error("resource uploadFile.Copy error(%v)", err)
		return
	}
	md5 := md5hash.Sum(nil)
	fMd5 := hex.EncodeToString(md5[:])
	file = &model.FileInfo{
		Md5:  fMd5,
		Type: fType,
		Size: int64(len(content)),
	}
	return
}

// Upload can upload a file object: store the info in Redis, and transfer the file to Bfs
func (s *Service) Upload(c context.Context, fileName string, fileType string, timing int64, body []byte) (location string, err error) {
	if location, err = s.dao.Upload(c, fileName, fileType, timing, body, s.c.Bfs); err != nil { // bfs
		log.Error("s.upload.UploadBfs() error(%v)", err)
	}
	return
}
