package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
)

const (
	_delete  = 1
	_maxRows = 20000
)

// WelfareTypeSave save welfare type
func (s *Service) WelfareTypeSave(id int, name, username string) (err error) {
	wt := new(model.WelfareType)
	wt.ID = id
	wt.Name = name
	wt.OperName = username

	if id == 0 {
		if err = s.dao.WelfareTypeAdd(wt); err != nil {
			log.Error("WelfareTypeAdd(%v) Error(%v)", wt, err)
		}
	} else {
		if err = s.dao.WelfareTypeUpd(wt); err != nil {
			log.Error("WelfareTypeUpd(%v) Error(%v)", wt, err)
		}
	}

	return
}

// WelfareTypeState delete welfare type
func (s *Service) WelfareTypeState(c context.Context, id int, username string) (err error) {
	tx := s.dao.BeginGormTran(c)

	if err = s.dao.WelfareTypeState(tx, id, _delete, 0, username); err != nil {
		log.Error("WelfareTypeState id(%v) Error(%v)", id, err)
		tx.Rollback()
		return
	}
	if err = s.dao.ResetWelfareTid(tx, id); err != nil {
		log.Error("ResetWelfareTid tid(%v) Error(%v)", id, err)
		tx.Rollback()
		return
	}

	return tx.Commit().Error
}

// WelfareTypeList get welfare type list
func (s *Service) WelfareTypeList() (wts []*model.WelfareTypeRes, err error) {
	if wts, err = s.dao.WelfareTypeList(); err != nil {
		log.Error("WelfareTypeList Error(%v)", err)
	}
	return
}

// WelfareSave save welfare
func (s *Service) WelfareSave(username string, req *model.WelfareReq) (err error) {
	var (
		burl string
		hurl string
	)
	wf := new(model.Welfare)
	copyFiled(wf, req)
	wf.OperName = username
	if burl, err = getRelativePath(wf.BackdropUri); err != nil {
		return
	}
	wf.BackdropUri = burl
	if hurl, err = getRelativePath(wf.HomepageUri); err != nil {
		return
	}
	wf.HomepageUri = hurl

	if req.ID == 0 {
		if err = s.dao.WelfareAdd(wf); err != nil {
			log.Error("WelfareAdd(%v) Error(%v)", wf, err)
		}
	} else {
		req.BackdropUri = burl
		req.HomepageUri = hurl
		if err = s.dao.WelfareUpd(req); err != nil {
			log.Error("WelfareUpd(%v) Error(%v)", wf, err)
		}
	}
	return
}

// WelfareState delete welfare
func (s *Service) WelfareState(id int, username string) (err error) {
	if err = s.dao.WelfareState(id, _delete, 0, username); err != nil {
		log.Error("WelfareState(%v) Error(%v)", id, err)
	}

	return
}

// WelfareList get welfare list
func (s *Service) WelfareList(tid int) (ws []*model.WelfareRes, err error) {
	if ws, err = s.dao.WelfareList(tid); err != nil {
		log.Error("WelfareList tid(%v) Error(%v)", tid, err)
		return
	}

	randomBFSHost := fmt.Sprintf(s.c.Property.WelfareBgHost, rand.Intn(3))
	for _, w := range ws {
		wbs, err := s.dao.WelfareBatchList(w.ID)
		if err != nil {
			log.Error("WelfareBatchList wid(%v) Error(%v)", w.ID, err)
			return ws, err
		}
		w.HomepageUri = fmt.Sprintf("%v%v", randomBFSHost, w.HomepageUri)
		w.BackdropUri = fmt.Sprintf("%v%v", randomBFSHost, w.BackdropUri)
		for _, wb := range wbs {
			w.ReceivedCount += wb.ReceivedCount
			w.Count += wb.Count
		}
	}

	return
}

// WelfareBatchUpload save upload welfare code
func (s *Service) WelfareBatchUpload(body []byte, name, username string, wid, vtime int) (err error) {
	wcb := new(model.WelfareCodeBatch)
	wcb.BatchName = name
	wcb.Wid = wid
	wcb.Vtime = time.Time(vtime)
	wcb.OperName = username

	wcs := make([]*model.WelfareCode, 0)

	str := string(body)
	for _, lineStr := range strings.Split(str, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		wc := new(model.WelfareCode)
		wc.Wid = wid
		wc.Code = lineStr
		wcs = append(wcs, wc)
	}

	wcb.Count = len(wcs)
	if wcb.Count > _maxRows {
		err = ecode.VipWelfareUploadMaxErr
		return
	}

	if err = s.dao.WelfareBatchSave(wcb); err != nil {
		log.Error("WelfareBatchSave (%v) Error(%v)", wcb, err)
		return
	}
	for _, wc := range wcs {
		wc.Bid = wcb.ID
	}
	if err = s.dao.WelfareCodeBatchInsert(wcs); err != nil {
		log.Error("WelfareBatchSave Error(%v)", err)
	}

	return
}

// WelfareBatchList get welfare batch list
func (s *Service) WelfareBatchList(wid int) (wbs []*model.WelfareBatchRes, err error) {
	if wbs, err = s.dao.WelfareBatchList(wid); err != nil {
		log.Error("WelfareBatchList wid(%v) Error(%v)", wid, err)
	}

	return
}

// WelfareBatchState delete welfare batch
func (s *Service) WelfareBatchState(c context.Context, id int, username string) (err error) {
	tx := s.dao.BeginGormTran(c)
	if err = s.dao.WelfareBatchState(tx, id, _delete, 0, username); err != nil {
		log.Error("WelfareBatchState(%v) Error(%v)", id, err)
		tx.Rollback()
		return
	}

	if err = s.dao.WelfareCodeStatus(tx, id, _delete); err != nil {
		log.Error("WelfareCodeStatus bid(%v) Error(%v)", id, err)
		tx.Rollback()
		return
	}

	return tx.Commit().Error
}

func copyFiled(wf *model.Welfare, req *model.WelfareReq) {
	wf.ID = req.ID
	wf.WelfareName = req.WelfareName
	wf.WelfareDesc = req.WelfareDesc
	wf.HomepageUri = req.HomepageUri
	wf.BackdropUri = req.BackdropUri
	wf.Recommend = req.Recommend
	wf.Rank = req.Rank
	wf.Tid = req.Tid
	wf.Stime = req.Stime
	wf.Etime = req.Etime
	wf.UsageForm = req.UsageForm
	wf.ReceiveRate = req.ReceiveRate
	wf.VipType = req.VipType
}

// getRelativePath get relative path
func getRelativePath(absolutePath string) (relativePath string, err error) {
	u, err := url.Parse(absolutePath)
	if err != nil {
		err = ecode.VipWelfareUrlUnvalid
		log.Error("hostChange ParseURL(%v) error (%v)", absolutePath, err)
		return
	}
	relativePath = u.Path
	return
}
