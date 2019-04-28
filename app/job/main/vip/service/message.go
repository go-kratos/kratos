package service

import (
	"context"
	"fmt"
	"strings"
	xtime "time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"
	"go-common/library/time"

	"github.com/pkg/errors"
)

const (
	maxSendMsgCount  = 200
	willExpiredDay   = 7
	hadExpiredDays   = 14
	willExpiredMsg   = "您还有%d天大会员即将到期，请尽快续期，享受更多特权！"
	willExpiredTitle = "大会员即将到期提醒"

	hadExpiredMsg   = "很抱歉的通知您，您的大会员已过期，请尽快续期享受更多特权！"
	hadExpiredTitle = "大会员过期提醒"

	vipWillExpiredMsgCode = "10_1_2"
	vipHadExpiredMsgCode  = "10_1_3"

	systemNotify = 4
)

func (s *Service) hadExpiredMsgJob() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.hadExpiredMsgJob panic(%v)", x)
			go s.hadExpiredMsgJob()
			log.Info("service.hadExpiredMsgJob recover")
		}
	}()
	log.Info("start send had expire msg job ...................")
	var (
		err  error
		mids []int64
	)
	now := xtime.Now()
	startTime := now.AddDate(0, 0, -hadExpiredDays-1)
	endTime := now.AddDate(0, 0, -hadExpiredDays)

	if mids, err = s.willExpireUser(time.Time(startTime.Unix()), time.Time(endTime.Unix()), model.VipStatusOverTime); err != nil {
		log.Error("will expire user(startTime:%v endTime:%v status:%v) error(%+v)", startTime, endTime, model.VipStatusOverTime, err)
		return
	}
	log.Info("send startTime(%v) endDate(%v) mids(%v)", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), mids)
	if err = s.batchSendMsg(mids, hadExpiredMsg, hadExpiredTitle, vipHadExpiredMsgCode, systemNotify); err != nil {
		log.Error("batch send msg error(%+v)", err)
		return
	}
	log.Info("end send had expire msg job..........................")
}

func (s *Service) willExpiredMsgJob() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.hadExpiredMsgJob panic(%v)", x)
			go s.willExpiredMsgJob()
			log.Info("service.hadExpiredMsgJob recover")
		}
	}()
	log.Info("start send will  expire msg job............................")
	var (
		err  error
		mids []int64
	)
	now := xtime.Now()
	startTime := now.AddDate(0, 0, willExpiredDay)
	endTime := now.AddDate(0, 0, willExpiredDay+1)
	if mids, err = s.willExpireUser(time.Time(startTime.Unix()), time.Time(endTime.Unix()), model.VipStatusNotOverTime); err != nil {
		log.Error("will expire user(startTime:%v endTime:%v status:%v) error(%+v)", startTime, endTime, model.VipStatusNotOverTime, err)
		return
	}
	log.Info("send startTime(%v) endDate(%v) mids(%v)", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"), mids)
	if err = s.batchSendMsg(mids, fmt.Sprintf(willExpiredMsg, willExpiredDay), willExpiredTitle, vipWillExpiredMsgCode, systemNotify); err != nil {
		log.Error("batch send msg error(%+v)", err)
		return
	}
	log.Info("end send will  expire msg job............................")
}

func (s *Service) willExpireUser(startTime time.Time, endTime time.Time, status int) (mids []int64, err error) {
	var (
		maxID int
		size  = 10000
	)

	if maxID, err = s.dao.SelMaxID(context.TODO()); err != nil {
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
		var tempMid []int64
		if tempMid, err = s.dao.SelVipUserInfos(context.TODO(), startID, endID, startTime, endTime, status); err != nil {
			err = errors.WithStack(err)
			return
		}
		mids = append(mids, tempMid...)
	}
	return
}

func (s *Service) batchSendMsg(mids []int64, content string, title string, ms string, dataType int) (err error) {
	if len(mids) <= maxSendMsgCount && len(mids) >= 1 {
		var midsStr = ""
		for _, v := range mids {
			midsStr += fmt.Sprintf(",%v", v)
		}
		if err = s.dao.SendMultipMsg(context.TODO(), midsStr, content, title, ms, dataType); err != nil {
			err = errors.WithStack(err)
			return
		}
	} else if len(mids) > maxSendMsgCount {
		page := len(mids) / maxSendMsgCount
		if len(mids)%maxSendMsgCount != 0 {
			page++
		}
		for i := 0; i < page; i++ {
			start := i * maxSendMsgCount
			end := (i + 1) * maxSendMsgCount
			if len(mids) < end {
				end = len(mids)
			}
			tempMids := mids[start:end]

			var midsStr []string
			for _, v := range tempMids {
				midsStr = append(midsStr, fmt.Sprintf("%v", v))
			}
			if err = s.dao.SendMultipMsg(context.TODO(), strings.Join(midsStr, ","), content, title, ms, dataType); err != nil {
				err = errors.WithStack(err)
				continue
			}
		}
	}
	return
}
