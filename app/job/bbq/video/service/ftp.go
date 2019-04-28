package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"go-common/app/job/bbq/video/model"
	"go-common/library/log"

	ftp "github.com/ftp-master"
)

const (
	filePathKey      = "search"
	videoFileName    = "bbqvideo"
	videoFileMD5Name = "bbqvideo.md5"
	userFileName     = "bbquser"
	userFileMD5Name  = "bbquser.md5"
	sugFileName      = "bbqsug"
	sugFileMD5Name   = "bbqsug.md5"
	sugSrcIDIdx      = 0
	sugSrcTermIdx    = 1
	sugSrcTypeIdx    = 2
	sugSrcScoreIdx   = 3
	defaultSugPath   = "/src/go-common/app/job/bbq/video/cmd/sug"
)

// writeSug .
func (s *Service) writeSug(path string, filename string, md5name string) (err error) {
	var num int64
	filePath := path + filename
	file, err := os.OpenFile(filePath, os.O_CREATE|syscall.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("writeSug os.OpenFile(%s) error(%v)", filename, err)
		return
	}
	defer file.Close()
	srcPath := s.c.FTP.LocalPath["sugsrc"]
	if srcPath == "" {
		srcPath = os.Getenv("GOPATH") + defaultSugPath
	}
	if srcPath == "" {
		log.Error("sugsrc path is empty")
		return
	}
	src, err := os.Open(srcPath)
	if err != nil {
		log.Error("writeSug os.Open source sug error(%v)", err)
		return
	}
	defer src.Close()
	br := bufio.NewReader(src)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		strArr := strings.Split(string(a), ",")
		var bs []byte
		var b []byte
		id, _ := strconv.ParseInt(strArr[sugSrcIDIdx], 10, 64)
		score, _ := strconv.ParseInt(strArr[sugSrcScoreIdx], 10, 64)
		sd := &model.Sug{
			ID:    id,
			Term:  strArr[sugSrcTermIdx],
			Type:  strArr[sugSrcTypeIdx],
			Score: score,
		}
		if b, err = json.Marshal(sd); err != nil {
			log.Warn("json.Marshal(%v) error(%v)", sd, err)
			continue
		}
		bs = append(bs, b...)
		bs = append(bs, '\n')
		if _, err = file.Write(bs); err != nil {
			log.Error("writeSug file.Write error(%v)", err)
			return
		}
		num++
		// log.Info("write sug term:[%s]", strArr[sugSrcTermIdx])
	}
	log.Info("writeSug success num (%d)", num)
	if err = s.writeMD5(path, filename, md5name); err != nil {
		return
	}
	s.uploadFile(path, filename, md5name)
	return
}

// writeUserInfo .
func (s *Service) writeUserInfo(path string, filename string, md5name string) (err error) {
	var lastID int64
	var num int64
	filePath := path + filename
	file, err := os.OpenFile(filePath, os.O_CREATE|syscall.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("writeUserInfo os.OpenFile(%s) error(%v)", filename, err)
		return
	}
	defer file.Close()
	for {
		var (
			dataList []*model.UserBaseDB
			bs       []byte
		)
		if dataList, err = s.dao.UsersByLast(context.Background(), lastID); err != nil {
			return
		}
		if len(dataList) == 0 {
			break
		}
		for _, data := range dataList {
			var b []byte
			sd := &model.UserSearch{
				ID:    data.MID,
				Uname: data.Uname,
			}
			if b, err = json.Marshal(sd); err != nil {
				log.Warn("json.Marshal(%v) error(%v)", sd, err)
				continue
			}
			bs = append(bs, b...)
			bs = append(bs, '\n')
			lastID = data.ID
			num++
			// log.Info("append user mid:[%d]", data.MID)
		}
		if _, err = file.Write(bs); err != nil {
			log.Error("writeUserInfo file.Write error(%v)", err)
			return
		}
	}
	log.Info("writeUserInfo success num (%d)", num)
	if err = s.writeMD5(path, filename, md5name); err != nil {
		return
	}
	s.uploadFile(path, filename, md5name)
	return
}

// writeVideoInfo .
func (s *Service) writeVideoInfo(path string) (err error) {
	var lastID int64
	var num int64
	var RelateID int64
	filePath := path + videoFileName
	file, err := os.OpenFile(filePath, os.O_CREATE|syscall.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("writeVideoInfo os.OpenFile(%s) error(%v)", videoFileName, err)
		return
	}
	defer file.Close()
	for {
		var (
			dataList []*model.VideoDB
			bs       []byte
		)
		if dataList, err = s.dao.VideosByLast(context.Background(), lastID); err != nil {
			return
		}
		if len(dataList) == 0 {
			break
		}
		for _, data := range dataList {
			var (
				b []byte
			)
			RelateID, err = strconv.ParseInt(fmt.Sprintf("%d00", data.AutoID), 10, 32)
			if err != nil {
				log.Error("relate id parse err [%v]", err)
				continue
			}
			sd := &model.VideoSearch{
				ID:      int32(RelateID),
				Title:   data.Title,
				PubTime: int32(data.Pubtime),
			}
			if b, err = json.Marshal(sd); err != nil {
				log.Warn("json.Marshal(%v) error(%v)", sd, err)
				continue
			}
			bs = append(bs, b...)
			bs = append(bs, '\n')
			lastID = data.AutoID
			num++
			// log.Info("append video svid:[%d]", data.ID)
		}
		if _, err = file.Write(bs); err != nil {
			log.Error("writeVideoInfo file.Write error(%v)", err)
			return
		}
	}
	log.Info("writeVideoInfo success num (%d)", num)
	if err = s.writeMD5(path, videoFileName, videoFileMD5Name); err != nil {
		return
	}
	s.uploadFile(path, videoFileName, videoFileMD5Name)
	return
}

// writeMD5 .
func (s *Service) writeMD5(path string, fname string, md5name string) (err error) {
	var out bytes.Buffer
	f, err := exec.LookPath("md5sum")
	if err != nil {
		log.Error("writeMD5 exec.LookPath(md5sum) error(%v)", err)
		f = "md5sum"
	}
	cmd := exec.Command(f, path+fname)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Error("writeMD5 cmd.Run() error(%v)", err)
		return
	}
	outStrs := strings.Split(out.String(), " ")
	if len(outStrs) == 0 {
		log.Error("writeMD5 len(outStrs) == 0")
		return errors.New("NONE MD5")
	}
	md5FileName := path + md5name
	md5File, err := os.OpenFile(md5FileName, os.O_CREATE|syscall.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Error("writeMD5 os.OpenFile(md5 file) error(%v)", err)
		return
	}
	defer md5File.Close()
	_, err = md5File.WriteString(outStrs[0])
	return
}

// SyncVideo2Search 同步视频
func (s *Service) SyncVideo2Search() {
	err := s.writeVideoInfo(s.c.FTP.LocalPath[filePathKey])
	if err != nil {
		log.Error("SyncVideo err[%v]", err)
	} else {
		log.Info("SyncVideo complete")
	}
}

// SyncUser2Search 同步用户
func (s *Service) SyncUser2Search() {
	err := s.writeUserInfo(s.c.FTP.LocalPath[filePathKey], userFileName, userFileMD5Name)
	if err != nil {
		log.Error("syncUser err[%v]", err)
	} else {
		log.Info("syncUser complete")
	}
}

// SyncSug2Search 同步sug
func (s *Service) SyncSug2Search() {
	err := s.writeSug(s.c.FTP.LocalPath[filePathKey], sugFileName, sugFileMD5Name)
	if err != nil {
		log.Error("syncSug err[%v]", err)
	} else {
		log.Info("syncSug complete")
	}
}

// SyncSearch 全量同步搜索各文件
func (s *Service) SyncSearch() {
	for {
		s.searchChan <- "SyncSearch Complete"
		var wg sync.WaitGroup
		fmt.Println("SyncSearch Start")
		stratAt := time.Now()
		wg.Add(3)
		go func() {
			defer wg.Done()
			s.SyncVideo2Search()
		}()
		go func() {
			defer wg.Done()
			s.SyncUser2Search()
		}()
		go func() {
			defer wg.Done()
			s.SyncSug2Search()
		}()
		wg.Wait()
		//耗时统计
		elapsed := time.Since(stratAt)
		log.Info("SyncSearch elapsed: %s", elapsed)
		fmt.Println(<-s.searchChan)
		//间隔1min
		time.Sleep(time.Minute * 1)
	}
}

// uploadFile .
func (s *Service) uploadFile(path string, fname string, md5name string) (err error) {
	ftp, err := ftp.Connect(s.c.FTP.Addr)
	if err != nil {
		log.Error("connect to ftp(%s) error(%v)", s.c.FTP.Addr, err)
		return
	}
	defer ftp.Quit()
	err = ftp.Login(s.c.FTP.User, s.c.FTP.Password)
	if err != nil {
		log.Error("ftp login(user:%s) error(%v)", s.c.FTP.User, err)
		return
	}
	defer ftp.Logout()
	ftpDir := fmt.Sprintf(s.c.FTP.RemotePath[filePathKey], fname)
	err = ftp.ChangeDir(ftpDir)
	if err != nil {
		log.Error("enter ftp path [%s] err [%v]", ftpDir, err)
		return
	}
	// delete source file, and upload file to ftp.
	if err = ftp.Delete(fname); err != nil {
		log.Error("ftp.Delete(%s) error(%v)", fname, err)
	}
	filePath := path + fname
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("os.Open(%s) error(%v)", filePath, err)
		return
	}
	defer file.Close()
	if err = ftp.Stor(fname, file); err != nil {
		log.Error("ftp.Stor(%s) error(%v)", fname, err)
		return
	}
	log.Info("ftp upload success(%s)", fname)
	// delete source.md5 file, and upload file to ftp.
	if err = ftp.Delete(md5name); err != nil {
		log.Error("ftp.Delete(%s) error(%v)", md5name, err)
	}
	md5FilePath := path + md5name
	md5File, err := os.Open(md5FilePath)
	if err != nil {
		log.Error("os.Open(%s) error(%v)", md5name, err)
		return
	}
	defer md5File.Close()
	if err = ftp.Stor(md5name, md5File); err != nil {
		log.Error("ftp.Stor(%s) error(%v)", md5name, err)
		return
	}
	log.Info("ftp upload success(%s)", md5name)
	return
}
