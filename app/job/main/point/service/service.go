package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-common/app/job/main/point/conf"
	"go-common/app/job/main/point/dao"
	"go-common/app/job/main/point/model"
	rpcmdl "go-common/app/service/main/point/model"
	"go-common/app/service/main/point/rpc/client"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

const (
	_vipPointChange = "vip_point_change_history"
	_pointChange    = "point_change_history"
	_insert         = "insert"
)

// Service struct
type Service struct {
	c                  *conf.Config
	dao                *dao.Dao
	oldVipPointDatabus *databus.Databus
	pointDatabus       *databus.Databus
	pointUpdate        *databus.Databus
	waiter             sync.WaitGroup
	closed             bool
	pointRPC           *client.Service
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		pointRPC: client.New(c.PointRPC),
	}
	if c.DataBus.OldVipBinlog != nil {
		s.oldVipPointDatabus = databus.New(c.DataBus.OldVipBinlog)
		s.waiter.Add(1)
		go s.syncoldpointdataproc()
	}
	if c.DataBus.PointBinlog != nil {
		s.pointDatabus = databus.New(c.DataBus.PointBinlog)
		s.waiter.Add(1)
		go s.pointbinlogproc()
	}
	if c.DataBus.PointUpdate != nil {
		s.pointUpdate = databus.New(c.DataBus.PointUpdate)
		s.waiter.Add(1)
		go s.pointupdateproc()
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() (err error) {
	defer s.waiter.Wait()
	s.closed = true
	s.dao.Close()
	if err = s.oldVipPointDatabus.Close(); err != nil {
		log.Error("s.oldVipPointDatabus.Close() error(%v)", err)
		return
	}
	if err = s.pointDatabus.Close(); err != nil {
		log.Error("s.pointDatabus.Close() error(%v)", err)
		return
	}
	if err = s.pointUpdate.Close(); err != nil {
		log.Error("s.pointUpdate.Close() error(%v)", err)
		return
	}
	return
}

func (s *Service) syncoldpointdataproc() {
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.oldVipPointDatabus.Messages()
		ok      bool
		c       = context.Background()
	)
	for {
		msg, ok = <-msgChan
		if s.closed || !ok {
			log.Info("syncoldpointdataproc msgChan closed")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		message := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), message); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", message, err)
			continue
		}
		if message.Table == _vipPointChange {
			history := new(model.VipPointChangeHistoryMsg)
			if err = json.Unmarshal(message.New, history); err != nil {
				log.Error("json.Unmarshal val(%v) error(%v)", string(message.New), err)
				continue
			}
			if message.Action == _insert {
				var count int64
				if count, err = s.dao.HistoryCount(c, history.Mid, history.OrderID); err != nil {
					log.Error("update point(%v) history(%v)", err, history)
					continue
				}
				if count > 0 {
					log.Error("update point change history had repeat record(%v)", history)
					continue
				}
				changeTime, err := time.ParseInLocation("2006-01-02 15:04:05", history.ChangeTime, time.Local)
				if err != nil {
					log.Error("update point time.ParseInLocation(%s) error(%v)", history.ChangeTime, err)
					continue
				}
				ph := &model.PointHistory{
					Mid:          history.Mid,
					Point:        history.Point,
					OrderID:      history.OrderID,
					ChangeType:   int(history.ChangeType),
					ChangeTime:   xtime.Time(changeTime.Unix()),
					RelationID:   history.RelationID,
					PointBalance: history.PointBalance,
					Remark:       history.Remark,
					Operator:     history.Operator,
				}
				s.updatePointWithHistory(c, ph)
			}
		}
	}
}

func (s *Service) pointbinlogproc() {
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.pointDatabus.Messages()
		ok      bool
		c       = context.Background()
	)
	for {
		msg, ok = <-msgChan
		if s.closed || !ok {
			log.Info("pointbinlogproc msgChan closed")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		if v.Table != _pointChange || v.Action != _insert {
			continue
		}
		h := new(model.VipPointChangeHistoryMsg)
		if err = json.Unmarshal(v.New, h); err != nil {
			log.Error("json.Unmarshal val(%v) error(%v)", string(v.New), err)
			continue
		}
		log.Info("point change log %+v", h)
		s.Notify(c, h)
	}
}

func (s *Service) updatePointWithHistory(c context.Context, ph *model.PointHistory) (pointBalance int64, activePoint int64, err error) {
	var (
		tx *xsql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("update point mid(%d) %+v record(%v)", ph.Mid, err, ph)
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				log.Error("update point tx.Commit %+v record(%v)", err, ph)
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
	}()
	if pointBalance, err = s.updatePoint(c, tx, ph.Mid, ph.Point, ph); err != nil {
		log.Error("update point mid(%d) %+v record(%v)", ph.Mid, err, ph)
		return
	}
	if _, err = s.dao.InsertPointHistory(c, tx, ph); err != nil {
		log.Error("update point mid(%d) %+v record(%v)", ph.Mid, err, ph)
		return
	}
	return
}

func (s *Service) updatePoint(c context.Context, tx *xsql.Tx, mid, point int64, ph *model.PointHistory) (pb int64, err error) {
	var (
		pi  *model.PointInfo
		ver int64
		a   int64
	)
	if pi, err = s.dao.TxPointInfo(c, tx, mid); err != nil {
		log.Error("%+v", err)
		return
	}
	if point == 0 {
		return
	}
	if pi == nil {
		if point < 0 {
			point = 0
			log.Error("update point<0 mid(%d) (%v)", mid, ph)
		}
		pb = point
		pi = new(model.PointInfo)
		pi.Ver = 1
		pi.PointBalance = point
		pi.Mid = mid
		if a, err = s.dao.InsertPoint(c, tx, pi); err != nil {
			log.Error("%v", err)
			return
		}
		if a != 1 {
			err = fmt.Errorf("operation failed")
			return
		}

	} else {
		pb = pi.PointBalance + point
		if pb < 0 {
			pb = 0
			log.Error("update point<0 mid(%d)(%v)", mid, ph)
		}
		pi.PointBalance = pb
		ver = pi.Ver
		pi.Ver++
		if a, err = s.dao.UpdatePointInfo(c, tx, pi, ver); err != nil {
			log.Error("%v", err)
			return
		}
		if a != 1 {
			err = fmt.Errorf("operation failed")
			return
		}
	}
	return
}

func (s *Service) fixdata(mtimeStr string) (err error) {
	var (
		pis   []*model.PointInfo
		c     = context.TODO()
		mtime time.Time
	)
	log.Info("fixdata start ")
	if mtime, err = time.ParseInLocation("2006-01-02 15:04:05", mtimeStr, time.Local); err != nil {
		log.Error("fixdata ParseInLocation error(%v)", err)
		return
	}
	if pis, err = s.dao.MidsByMtime(c, mtime); err != nil {
		log.Error("fixdata err %+v ", err)
		return
	}
	log.Info("fixdata count %d ", len(pis))
	for _, pi := range pis {
		log.Info("fixdata ing %v ", pi)
		var point int64
		if point, err = s.dao.LastOneHistory(c, pi.Mid); err != nil {
			log.Error("fixdata history err %+v ", err)
			return
		}
		if point == pi.PointBalance {
			continue
		}
		log.Info("fixdata mid(%d) %+v ", pi.Mid, err)
		var af int64
		if af, err = s.dao.FixPointInfo(c, pi.Mid, point); err != nil {
			log.Error("fixdata FixPointInfo mid(%d) %+v ", pi.Mid, err)
			return
		}
		if af != 1 {
			log.Error("fixdata af!=1 mid(%d)", pi.Mid)
			return
		}
	}
	log.Info("fixdata end ")
	return
}

func (s *Service) pointupdateproc() {
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.pointUpdate.Messages()
		ok      bool
		c       = context.Background()
	)
	for {
		msg, ok = <-msgChan
		if !ok {
			log.Info("pointupdateproc msgChan closed")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		v := &rpcmdl.ArgPointAdd{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		s.pointRPC.PointAddByBp(c, v)
	}
}
