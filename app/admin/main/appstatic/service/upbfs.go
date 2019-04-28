package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"

	"go-common/app/admin/main/appstatic/model"
	"go-common/library/log"
)

// Upload can upload a file object: store the info in Redis, and transfer the file to Bfs
func (s *Service) Upload(c context.Context, fileName string, fileType string, timing int64, body []byte) (location string, err error) {
	if s.c.Cfg.Storage == "nas" { // nas
		if location, err = s.dao.UploadNas(c, fileName, body, s.c.Nas); err != nil {
			log.Error("s.upload.UploadNas() error(%v)", err)
		}
		return
	}
	if location, err = s.dao.Upload(c, fileName, fileType, timing, body, s.c.Bfs); err != nil { // bfs
		log.Error("s.upload.UploadBfs() error(%v)", err)
	}
	return
}

// AddFile inserts file info into DB and updates its resource version+1
func (s *Service) AddFile(c context.Context, file *model.ResourceFile, version int) (err error) {
	if err = s.DB.Create(file).Error; err != nil {
		log.Error("resSrv.DB.Create error(%v)", err)
		return
	}
	// the resource containing the file updates its version
	if err = s.DB.Model(&model.Resource{}).Where("id = ?", file.ResourceID).Update("version", version+1).Error; err != nil {
		log.Error("resSrv.Update Version error(%v)", err)
		return
	}
	return nil
}

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

// TypeCheck checks whether the file type is allowed
func (s *Service) TypeCheck(fType string) (canAllow bool) {
	allowed := s.c.Cfg.Filetypes
	if len(allowed) == 0 {
		return true
	}
	for _, v := range allowed {
		if v == fType {
			return true
		}
	}
	return false
}
