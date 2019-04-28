package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/job/main/coupon/conf"
	"go-common/app/job/main/coupon/dao"
	"go-common/app/job/main/coupon/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/robfig/cron"
)

const (
	_couponTable          = "coupon_info_"
	_orderTable           = "coupon_order"
	_couponAllowanceTable = "coupon_allowance_info"
	_updateAct            = "update"
	_insertAct            = "insert"
)

// Service struct
type Service struct {
	c             *conf.Config
	dao           *dao.Dao
	couponDatabus *databus.Databus
	waiter        sync.WaitGroup
	notifyChan    chan *model.NotifyParam
	close         bool
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		notifyChan: make(chan *model.NotifyParam, 10240),
	}
	if c.DataBus.CouponBinlog != nil {
		s.couponDatabus = databus.New(c.DataBus.CouponBinlog)
		s.waiter.Add(1)
		go s.couponbinlogproc()
	}
	go s.notifyproc()
	t := cron.New()
	t.AddFunc(s.c.Properties.CheckInUseCouponCron, s.CheckInUseCoupon)
	t.AddFunc(s.c.Properties.CheckInUseCouponCartoonCron, s.CheckOrderInPayCoupon)
	t.Start()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.close = true
	s.couponDatabus.Close()
	s.dao.Close()
	s.waiter.Wait()
}

func (s *Service) couponbinlogproc() {
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.couponDatabus.Messages()
		ok      bool
		c       = context.Background()
	)
	for {
		msg, ok = <-msgChan
		if !ok || s.close {
			log.Info("couponbinlogproc closed")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%+v)", err)
		}
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		log.Info("couponbinlogproc log(%+v)", v)
		if err = s.Notify(c, v); err != nil {
			log.Error("s.Notify(%v) err(%v)", v, err)
		}
	}
}

func (s *Service) notifyproc() {
	var (
		msg          *model.NotifyParam
		ticker       = time.NewTicker(time.Duration(s.c.Properties.NotifyTimeInterval))
		mergeMap     = make(map[string]*model.NotifyParam)
		maxMergeSize = 1000
		full         bool
		ok           bool
		err          error
	)
	for {
		select {
		case msg, ok = <-s.notifyChan:
			if !ok {
				log.Info("notifyproc msgChan closed")
				return
			}
			if msg == nil {
				continue
			}
			if _, ok := mergeMap[msg.CouponToken]; !ok {
				mergeMap[msg.CouponToken] = msg
			}
			if len(mergeMap) < maxMergeSize {
				continue
			}
			full = true
		case <-ticker.C:
		}
		if len(mergeMap) > 0 {
			for _, v := range mergeMap {
				log.Info("retry notify coupon arg(%v)", v)
				if err = s.CheckCouponDeliver(context.TODO(), v); err != nil {
					log.Error("CheckCouponDeliver fail arg(%v) err(%v)", v, err)
					v.NotifyCount++
					if v.NotifyCount < s.c.Properties.MaxRetries {
						s.notifyChan <- v
					}
				}
			}
			mergeMap = make(map[string]*model.NotifyParam)
		}
		if full {
			time.Sleep(time.Duration(s.c.Properties.NotifyTimeInterval))
		}
	}
}
