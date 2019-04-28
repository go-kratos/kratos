package service

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"go-common/app/job/main/member/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	csclice []chan int64

	maxmid   int64 = 310000000
	scanned  int64
	errCount int64
)

func (s *Service) makeChan(num int) {
	csclice = make([]chan int64, num)
	for i := 0; i < num; i++ {
		csclice[i] = make(chan int64, 10000)
	}
}

// dataCheckMids check mid
func (s *Service) dataCheckMids() {
	var (
		i int64
	)
	if s.c.SyncRange.End > maxmid {
		s.c.SyncRange.End = maxmid
	}
	if s.c.SyncRange.Start < 0 {
		s.c.SyncRange.Start = 0
	}

	for i = s.c.SyncRange.Start; i < s.c.SyncRange.End; i++ {
		csclice[i%30] <- i
	}
}

// dataFixer
func (s *Service) dataFixer(cs chan int64) {
	for {
		mids := make([]int64, 0, 10)
		for mid := range cs {
			mids = append(mids, mid)
			if len(mids) >= 5 {
				break
			}
			atomic.AddInt64(&scanned, 1)
		}
		s.fix(mids)
	}
}

func (s *Service) fix(mids []int64) {
	var (
		err  error
		accs = make(map[int64]*model.AccountInfo)
		errs = make(map[int64]map[string]bool)
		c    = context.TODO()
		base *model.BaseInfo
	)
	func() {
		defer func() {
			if r := recover(); r != nil {
				r = errors.WithStack(r.(error))
				log.Error("fixer: wocao jingran recover le error(%+v)", r)
				time.Sleep(10 * time.Second)
			}
			time.Sleep(10 * time.Millisecond)
		}()

		if accs, errs, err = s.dao.Accounts(c, mids); err != nil {
			log.Error("fixer: dao.AccountInfo mid(%v) res(%v) error(%v)", mids, accs, err)
			return
		}
		for mid, res := range accs {
			log.Error("fixer: mid(%d) res(%+v)", mid, res)
			if base, err = s.dao.BaseInfo(c, mid); err != nil {
				log.Error("fixer: s.dao.BaseInfo mid(%d) err(%v)", mid, err)
				continue
			}
			if base == nil {
				log.Error("fixer: dataCheckErr mid(%d) res(%v),base(%v),detail(%v)", mid, res, base)
				continue
			}
			// all fields are same
			if sameAccInfo(base, res) {
				log.Info("fixer: sameAccInfo mid(%d) result true continue", mid)
				continue
			}

			// increase errCount and logging
			bs, _ := json.Marshal(base)
			jres, _ := json.Marshal(res)
			atomic.AddInt64(&errCount, 1)
			log.Error("fixer: dataCheckFail mid(%d) base(%s),res(%s),errCount(%d)", mid, bs, jres, atomic.LoadInt64(&errCount))

			if _, ok := errs[mid]; !ok {
				log.Error("fixer,errs[%v] is not ok", mid)
				continue
			}

			if asoOK := errs[mid]["asoOK"]; asoOK && !sameName(base, res) && len(res.Name) > 0 {
				s.dao.SetName(c, mid, res.Name)
			}
		}
		log.Info("fixer: dataCheckRight mids(%v) scanned(%d) errCount(%d)", mids, scanned, atomic.LoadInt64(&errCount))
	}()
}
