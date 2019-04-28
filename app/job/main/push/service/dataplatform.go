package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	pamdl "go-common/app/admin/main/push/model"
	"go-common/app/job/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// txCond get a new condition by tx.
func (s *Service) txCond(oldStatus, newStatus int) (cond *pamdl.DPCondition, err error) {
	ctx := context.Background()
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTx(ctx); err != nil {
		log.Error("tx.BeginTx() error(%v)", err)
		return
	}
	if cond, err = s.dao.TxCondByStatus(tx, oldStatus); err != nil || cond == nil {
		if e := tx.Rollback(); e != nil {
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if err = s.dao.TxUpdateCondStatus(tx, cond.ID, newStatus); err != nil {
		if e := tx.Rollback(); e != nil {
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
	}
	return
}

// data platform query
func (s *Service) dpQueryproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		cond, err := s.txCond(pushmdl.DpCondStatusPrepared, pushmdl.DpCondStatusSubmitting)
		if err != nil || cond == nil {
			time.Sleep(time.Second)
			continue
		}
		for i := 0; i < _retry; i++ {
			if cond.StatusURL, err = s.dao.DpSubmitQuery(context.Background(), cond.SQL); err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			log.Error("data platform add query(%+v) error(%v)", cond, err)
			s.dao.UpdateDpCondStatus(context.Background(), cond.ID, pushmdl.DpCondStatusFailed)
			s.dao.UpdateTaskStatus(context.Background(), cond.Task, pushmdl.TaskStatusFailed)
			continue
		}
		cond.Status = pushmdl.DpCondStatusSubmitted
		for i := 0; i < _retry; i++ {
			if err = s.dao.UpdateDpCond(context.Background(), cond); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil {
			log.Error("data platform update condition(%+v) error(%v)", cond, err)
		}
		time.Sleep(time.Second)
	}
}

// data platform get file
func (s *Service) dpFileproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		cond, err := s.txCond(pushmdl.DpCondStatusSubmitted, pushmdl.DpCondStatusPolling)
		if err != nil || cond == nil {
			time.Sleep(time.Second)
			continue
		}
		var (
			path  string
			files []string
		)
		if files = s.dpCheckJob(cond); len(files) == 0 {
			continue
		}
		for i := 0; i < _retry; i++ {
			if path, err = s.dpDownloadFiles(cond, files); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil || path == "" {
			log.Error("data platform download query(%+v) file error(%v)", cond, err)
			s.dao.UpdateDpCondStatus(context.Background(), cond.ID, pushmdl.DpCondStatusFailed)
			s.dao.UpdateTaskStatus(context.Background(), cond.Task, pushmdl.TaskStatusFailed)
			continue
		}
		cond.File = path
		cond.Status = pushmdl.DpCondStatusDone
		for i := 0; i < _retry; i++ {
			if err = s.dao.UpdateDpCond(context.Background(), cond); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil {
			log.Error("data platform UpdateDpCond(%+v) error(%v)", cond, err)
			continue
		}
		for i := 0; i < _retry; i++ {
			if err = s.dao.UpdateTask(context.Background(), strconv.FormatInt(cond.Task, 10), path, pushmdl.TaskStatusPretreatmentPrepared); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil {
			log.Error("s.dao.UpdateTask(%d,%s,%d) error(%v)", cond.Task, path, pushmdl.TaskStatusPretreatmentPrepared)
		}
		time.Sleep(time.Second)
	}
}

func (s *Service) dpCheckJob(cond *pamdl.DPCondition) (files []string) {
	now := time.Now()
	for {
		if time.Since(now) > time.Duration(s.c.Job.DpPollingTime) {
			log.Error("polling stoped, more over than dpPollingTime, give job up")
			s.dao.UpdateDpCondStatus(context.Background(), cond.ID, pushmdl.DpCondStatusFailed)
			break
		}
		res, err := s.dao.DpCheckJob(context.Background(), cond.StatusURL)
		if err != nil {
			log.Error("s.dao.DpCheckJob(%s) error(%v)", cond.StatusURL, err)
			time.Sleep(time.Second)
			continue
		}
		if res.StatusID == model.CheckJobStatusDoing || res.StatusID == model.CheckJobStatusPending {
			log.Info("polling (%s) ing..., status(%d)", cond.StatusURL, res.StatusID)
			time.Sleep(5 * time.Second)
			continue
		}
		if res.StatusID == model.CheckJobStatusOk {
			if len(res.Files) == 0 {
				log.Info("polling (%s) success, no files found", cond.StatusURL)
				s.dao.UpdateDpCondStatus(context.Background(), cond.ID, pushmdl.DpCondStatusNoFile)
				break
			}
			files = res.Files
			log.Info("polling (%s) success, files(%d)", cond.StatusURL, len(files))
			return
		}
		if res.StatusID == model.CheckJobStatusErr {
			log.Error("polling (%s) error, res(%+v)", cond.StatusURL, res)
			s.dao.UpdateDpCondStatus(context.Background(), cond.ID, pushmdl.DpCondStatusFailed)
			break
		}
	}
	log.Error("polling cond(%d) error", cond.ID)
	s.dao.UpdateTaskStatus(context.Background(), cond.Task, pushmdl.TaskStatusFailed)
	return
}

func (s *Service) dpDownloadFiles(cond *pamdl.DPCondition, files []string) (path string, err error) {
	for i := 0; i < _retry; i++ {
		if err = s.dao.UpdateDpCondStatus(context.Background(), cond.ID, pushmdl.DpCondStatusDownloading); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		return
	}
	dir := fmt.Sprintf("%s/%s", strings.TrimSuffix(s.c.Job.MountDir, "/"), time.Now().Format("20060102"))
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			log.Error("os.IsNotExist(%s) error(%v)", dir, err)
			return
		}
		if err = os.MkdirAll(dir, 0777); err != nil {
			log.Error("os.MkdirAll(%s) error(%v)", dir, err)
			return
		}
	}
	name := strconv.FormatInt(time.Now().UnixNano(), 10)
	path = fmt.Sprintf("%s/%x", dir, md5.Sum([]byte(name)))
	for _, f := range files {
		if err = s.dpDownloadFile(f, path); err != nil {
			return
		}
	}
	return
}

func (s *Service) dpDownloadFile(url, path string) (err error) {
	var (
		res     []byte
		content [][]byte
	)
	for i := 0; i < _retry; i++ {
		if res, err = s.dao.DpDownloadFile(context.Background(), url); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Error("s.dao.DpDownloadFile(%s) error(%v)", url, err)
		return
	}
	for _, bs := range bytes.Split(res, []byte("\n")) {
		n := bytes.Split(bs, []byte("\u0001"))
		content = append(content, bytes.Join(n, []byte("	")))
	}
	for i := 0; i < _retry; i++ {
		if err = s.saveDpFile(path, bytes.Join(content, []byte("\n"))); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Error("s.saveNASFile(%s) error(%v)", url, err)
	}
	return
}

// saveDpFile writes data platform data into NAS.
func (s *Service) saveDpFile(path string, data []byte) (err error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("s.saveDpFile(%s) OpenFile() error(%v)", path, err)
		return
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		log.Error("s.saveDpFile(%s) f.Write() error(%v)", path, err)
	}
	return
}
