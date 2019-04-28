package pgc

import (
	"encoding/json"
	"os"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

const (
	//ContLimit is used for getting ugc value 50 records every time
	_ContLimit = 50
)

func (s *Service) seaPgcContproc() {
	for {
		if s.daoClosed {
			log.Info("seaPgcContproc DB closed!")
			return
		}
		s.seaPgcCont()
		time.Sleep(time.Duration(s.c.Search.Cfg.UploadFre))
	}
}

// seaPgcCont is used for generate search content content
func (s *Service) seaPgcCont() {
	var (
		err     error
		seasons []*model.SearPgcCon
		str     []byte // the json string to write in file
		cnt     int
		cycle   int //cycle count and every cycle is 50 records
		id      int
	)
	if cnt, err = s.dao.PgcContCount(ctx); err != nil {
		log.Error(errFormat, "searchCont", "OnlineSeasonsC", err)
		return
	}
	cycle = cnt / _ContLimit
	if cnt%_ContLimit != 0 {
		cycle = cnt/_ContLimit + 1
	}
	// write into the file
	file, error := os.OpenFile(s.c.Search.PgcContPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if error != nil {
		log.Error(errFormat, "searchSug", "OpenFile", err)
		return
	}
	//cycle get sql value
	for i := 0; i < cycle; i++ {
		if i == 0 {
			id = 0
		} else {
			id = seasons[len(seasons)-1].ID
		}
		if seasons, err = s.dao.PgcCont(ctx, id, _ContLimit); err != nil {
			log.Error(errFormat, "PgcCont", "PgcCont", err)
			return
		}
		for _, v := range seasons {
			if str, err = json.Marshal(v); err != nil {
				log.Error(errFormat, "searchSug", "JsonMarshal", err)
				return
			}
			file.WriteString(string(str) + "\n")
		}
	}
	file.Close()
	//calculate file's md5
	if err = s.ftpDao.FileMd5(s.c.Search.PgcContPath, s.c.Search.PgcContMd5Path); err != nil {
		log.Error(errFormat, "searPgcCont", "fileMd5", err)
		return
	}
	// upload original file
	if err = s.ftpDao.UploadFile(s.c.Search.PgcContPath, s.c.Search.FTP.RemotePgcCont, s.c.Search.FTP.RemotePgcURL); err != nil {
		log.Error(errFormat, "searPgcCont-File", "uploadFile", err)
		return
	}
	// upload md5 file
	if err = s.ftpDao.UploadFile(s.c.Search.PgcContMd5Path, s.c.Search.FTP.RemotePgcContMd5, s.c.Search.FTP.RemotePgcURL); err != nil {
		log.Error(errFormat, "searPgcCont-Md5", "uploadFile", err)
		return
	}
	log.Info("FTP Upload Success")
}
