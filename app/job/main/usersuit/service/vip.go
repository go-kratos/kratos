package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/usersuit/model"
	vipmdl "go-common/app/service/main/vip/model"
	"go-common/library/log"
)

const (
	_vipGid           = 31
	_vipUserInfoTable = "vip_user_info"
)

func (s *Service) vipconsumerproc() {
	defer s.wg.Done()
	var (
		msgs = s.vipBinLogSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.vipBinLogSub.Message closed")
			return
		}
		msg.Commit()
		m := &model.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		switch m.Table {
		case _vipUserInfoTable:
			if m.Action == "update" {
				s.dealUserPendantEquip(c, m.New, m.Old)
			}
		default:
			log.Warn("vipBinLogConsumer unknown message action(%s)", m.Table)
		}
		if err != nil {
			log.Error("vipBinLogMessage key(%s) value(%s) partition(%d) offset(%d) commit error(%v)", msg.Key, msg.Value, msg.Partition, msg.Offset, err)
			continue
		}
		log.Info("vipBinLogMessage key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
	}
}

func (s *Service) dealUserPendantEquip(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	mr := &model.VipInfoMessage{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	var (
		gid int64
		pe  *model.PendantEquip
	)
	if pe, err = s.pendantDao.PendantEquipMID(c, mr.Mid); err != nil {
		log.Error("mid(%d) s.pendantDao.PendantEquipMID error(%v)", mr.Mid, err)
		return
	}
	if pe == nil || pe.Pid == 0 || pe.Expires == 0 {
		log.Warn("mid(%d) no equip pendant(%d) expires(%d)", mr.Mid, pe.Pid, pe.Expires)
		return
	}
	if gid, err = s.pendantDao.PendantEquipGidPid(c, pe.Pid); err != nil {
		log.Error("mid(%d) pid(%d) s.pendantDao.PendantEquipGidPid error(%v)", mr.Mid, pe.Pid, err)
		return
	}
	if gid != _vipGid {
		log.Warn("mid(%d) no equip the vip gid(%d) of pid(%d)", mr.Mid, gid, pe.Pid)
		return
	}
	if mr.VipStatus == vipmdl.VipStatusNotOverTime {
		var ts time.Time
		if ts, err = time.ParseInLocation(model.TimeFormatSec, mr.VipOverdueTime, time.Local); err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", mr.VipOverdueTime, err)
			return
		}
		if ts.Unix() <= pe.Expires {
			log.Warn("mid(%d) pendant equip_time(%d) than vipoverdue_time(%d)", mr.Mid, pe.Expires, ts.Unix())
			return
		}
		if _, err = s.pendantDao.UpEquipExpires(c, mr.Mid, ts.Unix()); err != nil {
			log.Error("s.pendantDao.UpEquipExpires(%d,%d) error(%+v)", mr.Mid, ts.Unix(), err)
			return
		}
	} else {
		if _, err = s.pendantDao.UpEquipMID(c, mr.Mid); err != nil {
			log.Error("s.pendantDao.UpEquipMID(%d) error(%+v)", mr.Mid, err)
			return
		}
		log.Warn("mid(%d) vip status is overtime", mr.Mid)
	}
	s.pendantDao.DelEquipCache(c, mr.Mid)
	s.addNotify(func() {
		s.accNotify(context.TODO(), mr.Mid, model.AccountNotifyUpdatePendant)
	})
	return
}
