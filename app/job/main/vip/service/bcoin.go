package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

//HandlerBcoin handler bcoin
func (s *Service) HandlerBcoin() (err error) {
	var (
		batchSize int64 = 2000
		oldMaxID  int64
		newMaxID  int64
		exitMap   = make(map[string]int)
		batch     []*model.VipBcoinSalary
	)
	if oldMaxID, err = s.dao.SelOldBcoinMaxID(context.TODO()); err != nil {
		log.Error("s.dao.SelOldBcoinMaxID error(%+v)", err)
		return
	}
	if newMaxID, err = s.dao.SelBcoinMaxID(context.TODO()); err != nil {
		log.Error("s.dao.SelBcoinMaxID error(%v)", err)
		return
	}
	page := newMaxID / batchSize
	if newMaxID%batchSize != 0 {
		page++
	}
	for i := 0; i < int(page); i++ {
		arg := new(model.QueryBcoinSalary)
		arg.StartID = int64(i) * batchSize
		arg.EndID = int64(i+1) * batchSize
		if arg.EndID > newMaxID {
			arg.EndID = newMaxID
		}
		arg.GiveNowStatus = -1
		var res []*model.VipBcoinSalary
		if res, err = s.dao.SelBcoinSalaryData(context.TODO(), arg.StartID, arg.EndID); err != nil {
			log.Error("s.dao.SelBcoinSalary(%+v) error(%+v)", arg, err)
			return
		}
		for _, v := range res {
			exitMap[s.makeBcoinMD5(v)] = 1
		}
	}

	page = oldMaxID / batchSize
	if oldMaxID%batchSize != 0 {
		page++
	}
	for i := 0; i < int(page); i++ {
		startID := int64(i) * batchSize
		EndID := int64(i+1) * batchSize
		if EndID > oldMaxID {
			EndID = oldMaxID
		}
		var res []*model.VipBcoinSalary
		if res, err = s.dao.SelOldBcoinSalary(context.TODO(), startID, EndID); err != nil {
			log.Error("sel.OldBcoinSalary(startID:%v endID:%v) error(%+v)", startID, EndID, err)
			return
		}

		for _, v := range res {
			if exitMap[s.makeBcoinMD5(v)] == 0 {
				batch = append(batch, v)
			}
		}

		if err = s.dao.BatchAddBcoinSalary(batch); err != nil {
			log.Error("s.dao.BatchAddBcoinSalary (%+v)", err)
			return
		}
		batch = nil
	}
	return
}

func (s *Service) handleraddbcoinproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerupdaterechargeorderproc panic(%v)", x)
			go s.handleraddbcoinproc()
			log.Info("service.handlerupdaterechargeorderproc recover")
		}
	}()
	for {
		msg := <-s.handlerAddBcoinSalary
		log.Info("cur bcoin msage:%+v", msg)
		for i := 0; i < s.c.Property.Retry; i++ {
			if err := s.dao.AddBcoinSalary(context.TODO(), msg); err != nil {
				log.Error("s.dao.addbcoinsalary(%+v) error(%+v)", msg, err)
			} else if err == nil {
				break
			}
		}

	}
}

func (s *Service) handlerdelbcoinproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerdelbcoinproc panic(%v)", x)
			go s.handlerdelbcoinproc()
			log.Info("service.handlerdelbcoinproc recover")
		}
	}()
	var err error
	for {
		msg := <-s.handlerDelBcoinSalary
		for i := 0; i < s.c.Property.Retry; i++ {
			if err = s.dao.DelBcoinSalary(context.TODO(), msg.Payday, msg.Mid); err == nil {
				break
			}
			log.Error("s.dao.DelBcoinSalary(msg:%+v) error(%+v)", msg, err)
		}
	}
}

func (s *Service) handlerupdatebcoinproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.handlerupdaterechargeorderproc panic(%v)", x)
			go s.handlerupdatebcoinproc()
			log.Info("service.handlerupdaterechargeorderproc recover")
		}
	}()
	var err error
	for {
		msg := <-s.handlerUpdateBcoinSalary
		for i := 0; i < s.c.Property.Retry; i++ {
			if err = s.dao.UpdateBcoinSalary(context.TODO(), msg.Payday, msg.Mid, msg.Status); err == nil {
				break
			}
			log.Error("s.dao.UpdateBcoinSalary(msg:%+v) error(%+v)", msg, err)
		}
	}
}

func (s *Service) sendBcoinJob() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	log.Info("sen bcoin job start ........................................ ")
	s.sendBcoin()
	log.Info("sen bcoin job end ........................................ ")
}

func (s *Service) sendBcoin() {
	var (
		maxID     int64
		batchSize int64 = 3000
		sendSize        = 50
		err       error
	)
	if maxID, err = s.dao.SelBcoinMaxID(context.TODO()); err != nil {
		log.Error("s.dao.selBcoinMaxID() error(%+v)", err)
		return
	}
	page := maxID / batchSize
	if maxID%batchSize != 0 {
		page++
	}
	now := time.Now()
	startMonth := now.AddDate(0, 0, 1-now.Day())
	endMonth := startMonth.AddDate(0, 1, 0)
	sendInfo := s.sendInfo()
	for i := 0; i < int(page); i++ {
		arg := new(model.QueryBcoinSalary)
		arg.StartID = int64(i) * batchSize
		arg.EndID = int64(i+1) * batchSize
		arg.GiveNowStatus = 0
		arg.Status = 0
		arg.StartMonth = xtime.Time(startMonth.Unix())
		arg.EndMonth = xtime.Time(endMonth.Unix())
		var res []*model.VipBcoinSalary
		if res, err = s.dao.SelBcoinSalary(context.TODO(), arg); err != nil {
			log.Error("s.dao.selBcoinSalary(%+v) error(%+v)", arg, err)
			return
		}
		pageSend := len(res) / sendSize
		if len(res)%sendSize != 0 {
			pageSend++
		}
		for j := 0; j < pageSend; j++ {
			start := j * sendSize
			end := int(j+1) * sendSize
			if end > len(res) {
				end = len(res)
			}
			if err = s.sendBocinNow(res[start:end], sendInfo.Amount, sendInfo.DueDate); err != nil {
				log.Error("%+v", err)
				return
			}

		}

	}
}

func (s *Service) sendInfo() (r *model.BcoinSendInfo) {
	var (
		c      time.Time
		day    = s.c.Property.AnnualVipBcoinDay
		amount = s.c.Property.AnnualVipBcoinCouponMoney
	)
	r = new(model.BcoinSendInfo)
	r.Amount = int32(amount)
	r.DayOfMonth = day
	c = time.Now()
	c = c.AddDate(0, 1, int(day)-1-c.Day())
	r.DueDate = xtime.Time(c.Unix())
	return
}

func (s *Service) sendBocinNow(res []*model.VipBcoinSalary, amount int32, duTime xtime.Time) (err error) {
	var (
		mids []int64
		ids  []int64
	)
	for _, v := range res {
		mids = append(mids, v.Mid)
		ids = append(ids, v.ID)
	}
	if err = s.dao.SendBcoin(context.TODO(), mids, amount, duTime, "127.0.0.1"); err != nil {
		err = errors.WithStack(err)
		return
	}

	if err = s.dao.UpdateBcoinSalaryBatch(context.TODO(), ids, 1); err != nil {
		err = errors.WithStack(err)
		return
	}
	return

}

func (s *Service) makeBcoinMD5(r *model.VipBcoinSalary) string {
	key := fmt.Sprintf("%v,%v,%v,%v,%v,%v", r.Mid, r.Memo, r.Amount, r.Payday.Time().Format("2006-01-02"), r.GiveNowStatus, r.Status)
	hash := md5.New()
	hash.Write([]byte(key))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
