package service

import (
	"context"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_fail = 2

	_handlering = 2
	_finish     = 3

	_nomarl = 0

	_statusnomarl = 1
)

func (s *Service) pushDataJob() {
	log.Info("push data job start..................")
	if succeed := s.dao.AddTransferLock(context.TODO(), "lock:pushDatajob"); succeed {
		if err := s.pushData(context.TODO()); err != nil {
			log.Error("error(%+v)", err)
		}
	}
	log.Info("push data job end.....................")
}

func (s *Service) pushData(c context.Context) (err error) {
	var (
		res         []*model.VipPushData
		pushDataMap = make(map[int64]*model.VipPushData)
		pushMidsMap = make(map[int64][]int64)
		maxID       int
		size        = s.c.Property.BatchSize
		vips        []*model.VipUserInfo
		curDate     time.Time
		rel         *model.VipPushResq
	)
	now := time.Now()
	format := now.Format("2006-01-02")
	if curDate, err = time.ParseInLocation("2006-01-02", format, time.Local); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res, err = s.dao.PushDatas(c, format); err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(res) == 0 {
		log.Info("not need reduce push data.........")
		return
	}
	for _, v := range res {
		pushDataMap[v.ID] = v
	}
	if maxID, err = s.dao.SelMaxID(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	page := maxID / size
	if maxID%size != 0 {
		page++
	}
	for i := 0; i < page; i++ {
		startID := i * size
		endID := (i + 1) * size
		if vips, err = s.dao.SelUserInfos(context.TODO(), startID, endID); err != nil {
			err = errors.WithStack(err)
			return
		}
		for _, v := range vips {
			for key, val := range pushDataMap {
				startDate := curDate.AddDate(0, 0, int(val.ExpiredDayStart))
				endDate := curDate.AddDate(0, 0, int(val.ExpiredDayEnd))
				if !(v.OverdueTime.Time().Before(startDate) || v.OverdueTime.Time().After(endDate)) && v.PayType == model.Normal && val.DisableType == _nomarl && val.Status != _fail {
					mids := pushMidsMap[key]
					mids = append(mids, v.Mid)
					pushMidsMap[key] = mids
				}
			}
		}
	}
	for key, val := range pushMidsMap {
		data := pushDataMap[key]
		var status int8
		progressStatus := data.ProgressStatus
		pushedCount := data.PushedCount
		if rel, err = s.dao.PushData(context.TODO(), val, data, format); err != nil {
			log.Error("push data error(%+v)", err)
			continue
		}
		if rel.Code != 0 {
			status = _fail
		} else {
			pushedCount++
			if pushedCount == data.PushTotalCount {
				progressStatus = _finish
			} else {
				progressStatus = _handlering
			}
			status = _statusnomarl
		}
		if err = s.dao.UpdatePushData(context.TODO(), status, progressStatus, pushedCount, rel.Code, rel.Data, data.ID); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}
