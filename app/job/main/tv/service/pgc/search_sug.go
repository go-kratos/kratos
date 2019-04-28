package pgc

import (
	"encoding/json"
	"os"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

const (
	errFormat = "Func:[%s] - Step:[%s] - Error:[%v]"
)

func (s *Service) pgcSeaSug(f *os.File) (err error) {
	var (
		str []byte // the json string to write in file
		sug []*model.SearchSug
	)
	if sug, err = s.dao.PgcSeaSug(ctx); err != nil {
		log.Error(errFormat, "searchSug", "PgcSeaSug", err)
		return
	}
	for _, v := range sug {
		if str, err = json.Marshal(v); err != nil {
			log.Error(errFormat, "searchSug", "JsonMarshal", err)
			return
		}
		f.WriteString(string(str) + "\n")
	}
	return
}

func (s *Service) ugcSeaSug(f *os.File) (err error) {
	var (
		str []byte // the json string to write in file
		sug []*model.SearchSug
	)
	if sug, err = s.dao.UgcSeaSug(ctx); err != nil {
		log.Error(errFormat, "ugcSeaSug", "UgcSeaSug", err)
		return
	}
	for _, v := range sug {
		if str, err = json.Marshal(v); err != nil {
			log.Error(errFormat, "ugcSeaSug", "JsonMarshal", err)
			return
		}
		f.WriteString(string(str) + "\n")
	}
	return
}

func (s *Service) searchSugproc() {
	for {
		if s.daoClosed {
			log.Info("searchSugproc DB closed!")
			return
		}
		s.searchSug()
		time.Sleep(time.Duration(s.c.Search.Cfg.UploadFre))
	}
}

// generate the valid seasons file for search suggestion
func (s *Service) searchSug() {
	// write into the file
	file, err := os.OpenFile(s.c.Search.SugPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		log.Error(errFormat, "searchSug", "OpenFile", err)
		return
	}
	if err := s.pgcSeaSug(file); err != nil {
		log.Error(errFormat, "searchSug", "OpenFile", err)
		return
	}
	//if switch is on then onpen ugc search suggest
	if s.c.Search.UgcSwitch == "on" {
		if err := s.ugcSeaSug(file); err != nil {
			log.Error(errFormat, "searchSug", "OpenFile", err)
			return
		}
	}
	file.Close()
	// calculate file's md5
	if err := s.ftpDao.FileMd5(s.c.Search.SugPath, s.c.Search.Md5Path); err != nil {
		log.Error(errFormat, "searchSug", "fileMd5", err)
		return
	}
	// upload original file
	if err := s.ftpDao.UploadFile(s.c.Search.SugPath, s.c.Search.FTP.RemoteFName, s.c.Search.FTP.URL); err != nil {
		log.Error(errFormat, "searchSug-File", "uploadFile", err)
		return
	}
	//upload md5 file
	if err := s.ftpDao.UploadFile(s.c.Search.Md5Path, s.c.Search.FTP.RemoteMd5, s.c.Search.FTP.URL); err != nil {
		log.Error(errFormat, "searchSug-Md5", "uploadFile", err)
		return
	}
	log.Error("FTP Upload Success")
}
