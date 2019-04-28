package ugc

import (
	"encoding/json"
	"os"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

const (
	errFormat = "Func:[%s] - Step:[%s] - Error:[%v]"
	//ContLimit is used for getting pgc season value 50 records every time
	_ContLimit = 50
)

func (s *Service) seaUgcContproc() {
	for {
		if s.daoClosed {
			log.Info("seaUgcContproc DB closed!")
			return
		}
		s.seaUgcCont()
		time.Sleep(time.Duration(s.c.Search.Cfg.UploadFre))
	}
}

// seaUgcCont is used for generate search content content
func (s *Service) seaUgcCont() {
	var (
		archives []*model.SearUgcCon
		str      []byte // the json string to write in file
		cycle    int    //cycle count and every cycle is 50 records
		maxID    = 0
		cnt      int
		file     *os.File
		err      error
	)
	if cnt, err = s.dao.UgcCnt(ctx); err != nil {
		log.Error(errFormat, "searUgcCont", "UgcCnt", err)
		return
	}
	cycle = cnt / _ContLimit
	if cnt%_ContLimit != 0 {
		cycle = cnt/_ContLimit + 1
	}
	// write into the file
	if file, err = os.OpenFile(s.c.Search.UgcContPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766); err != nil {
		log.Error(errFormat, "searUgcCont", "OpenFile", err)
		return
	}
	//cycle get sql value
	for i := 0; i < cycle; i++ {
		if archives, maxID, err = s.dao.UgcCont(ctx, maxID, _ContLimit); err != nil {
			log.Error(errFormat, "searUgcCont", "UgcCont", err)
			return
		}
		for _, v := range archives {
			v.Typeid = s.getPType(v.Typeid)
			if str, err = json.Marshal(v); err != nil {
				log.Error(errFormat, "searUgcCont", "JsonMarshal", err)
				return
			}
			file.WriteString(string(str) + "\n")
		}
	}
	file.Close()
	//calculate file's md5
	if err = s.ftpDao.FileMd5(s.c.Search.UgcContPath, s.c.Search.UgcContMd5Path); err != nil {
		log.Error(errFormat, "searUgcCont", "fileMd5", err)
		return
	}
	// upload original file
	if err = s.ftpDao.UploadFile(s.c.Search.UgcContPath, s.c.Search.FTP.RemoteUgcCont, s.c.Search.FTP.RemoteUgcURL); err != nil {
		log.Error(errFormat, "searUgcCont-File", "uploadFile", err)
		return
	}
	// upload md5 file
	if err = s.ftpDao.UploadFile(s.c.Search.UgcContMd5Path, s.c.Search.FTP.RemoteUgcContMd5, s.c.Search.FTP.RemoteUgcURL); err != nil {
		log.Error(errFormat, "searUgcCont-Md5", "uploadFile", err)
		return
	}
	log.Info("FTP Upload Success")
}
