package service

import (
	"context"

	"go-common/app/job/main/passport-user/model"
	"go-common/library/log"
	"time"
)

func (s *Service) getAsoAccount(start, end, count int64) {
	chanNum := int64(s.c.FullSync.AsoAccount.ChanNum)
	for {
		log.Info("getAsoAccount, start %d, end %d, count %d", start, end, count)
		var (
			res []*model.OriginAccount
			err error
		)
		for {
			if res, err = s.d.AsoAccount(context.Background(), start, count); err != nil {
				log.Error("fail to get AsoAccount error(%+v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
		for _, a := range res {
			s.asoAccountChan[a.Mid%chanNum] <- a
		}
		if start > end || len(res) == 0 {
			log.Info("sync asoAccount finished! endID(%d)", start)
			break
		}
		start = res[len(res)-1].Mid
	}
}

func (s *Service) asoAccountConsume(c chan *model.OriginAccount) {
	for {
		a, ok := <-c
		if !ok {
			log.Error("asoAccountChan closed")
			return
		}
		for i := 0; i < _retry; i++ {
			if err := s.syncAsoAccount(a); err != nil {
				log.Error("fail to sync asoAccount(%+v) error(%+v)", a, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}

func (s *Service) getAsoAccountInfo(start, end, count int64) {
	chanNum := int64(s.c.FullSync.AsoAccountInfo.ChanNum)
	initStart, initEnd, initCount := start, end, count
	for i := 0; i < _asoAccountInfoSharding; i++ {
		start, end, count = initStart, initEnd, initCount
		for {
			log.Info("getAsoAccountInfo, start %d, end %d, count %d, suffix %d", start, end, count, i)
			var (
				res []*model.OriginAccountInfo
				err error
			)
			for {
				if res, err = s.d.AsoAccountInfo(context.Background(), start, count, int64(i)); err != nil {
					log.Error("fail to get AsoAccountInfo error(%+v)", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}
			for _, a := range res {
				s.asoAccountInfoChan[a.ID%chanNum] <- a
			}
			if start > end || len(res) == 0 {
				log.Info("sync asoAccountInfo(%d) finished! endID(%d)", i, start)
				break
			}
			start = res[len(res)-1].ID
		}
	}
}

func (s *Service) asoAccountInfoConsume(c chan *model.OriginAccountInfo) {
	for {
		a, ok := <-c
		if !ok {
			log.Error("asoAccountInfoChan closed")
			return
		}
		for {
			if err := s.syncAsoAccountInfo(a); err != nil {
				log.Error("fail to sync asoAccountInfo(%+v) error(%+v)", a, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}

func (s *Service) getAsoAccountReg(start, end, count int64) {
	chanNum := int64(s.c.FullSync.AsoAccountReg.ChanNum)
	initStart, initEnd, initCount := start, end, count
	for i := 0; i < _asoAccountRegOriginSharding; i++ {
		start, end, count = initStart, initEnd, initCount
		for {
			log.Info("getAsoAccountReg, start %d, end %d, count %d, suffix %d", start, end, count, i)
			var (
				res []*model.OriginAccountReg
				err error
			)
			for {
				if res, err = s.d.AsoAccountReg(context.Background(), start, count, int64(i)); err != nil {
					log.Error("fail to get AsoAccountReg error(%+v)", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}
			for _, a := range res {
				s.asoAccountRegChan[a.ID%chanNum] <- a
			}
			if start > end || len(res) == 0 {
				log.Info("sync asoAccountReg(%d) finished! endID(%d)", i, start)
				break
			}
			start = res[len(res)-1].ID
		}
	}
}

func (s *Service) asoAccountRegConsume(c chan *model.OriginAccountReg) {
	for {
		a, ok := <-c
		if !ok {
			log.Error("asoAccountRegChan closed")
			return
		}
		for {
			if err := s.syncAsoAccountReg(a); err != nil {
				log.Error("fail to sync asoAccountReg(%+v) error(%+v)", a, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}

func (s *Service) getAsoAccountSns(start, end, count int64) {
	chanNum := int64(s.c.FullSync.AsoAccountReg.ChanNum)
	for {
		log.Info("getAsoAccountSns, start %d, end %d, count %d", start, end, count)
		var (
			res []*model.OriginAccountSns
			err error
		)
		for {
			if res, err = s.d.AsoAccountSns(context.Background(), start, count); err != nil {
				log.Error("fail to get AsoAccountSns error(%+v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
		for _, a := range res {
			s.asoAccountSnsChan[a.Mid%chanNum] <- a
		}
		if start > end || len(res) == 0 {
			log.Info("sync asoAccountSns finished! endID(%d)", start)
			break
		}
		start = res[len(res)-1].Mid
	}
}

func (s *Service) asoAccountSnsConsume(c chan *model.OriginAccountSns) {
	for {
		a, ok := <-c
		if !ok {
			log.Error("asoAccountSnsChan closed")
			return
		}
		for {
			if err := s.syncAsoAccountSns(a); err != nil {
				log.Error("fail to sync asoAccountSns(%+v) error(%+v)", a, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}

func (s *Service) getAsoTelBindLog(start, end, count int64) {
	chanNum := int64(s.c.FullSync.AsoAccountReg.ChanNum)
	for {
		log.Info("getAsoTelBindLog, start %d, end %d, count %d", start, end, count)
		var (
			res []*model.UserTel
			err error
		)
		for {
			if res, err = s.d.UserTel(context.Background(), start, count); err != nil {
				log.Error("fail to get UserTel error(%+v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
		for _, a := range res {
			s.asoTelBindLogChan[a.Mid%chanNum] <- a
		}
		if start > end || len(res) == 0 {
			log.Info("sync asoTelBindLog finished! endID(%d)", start)
			break
		}
		start = res[len(res)-1].Mid
	}
}

func (s *Service) asoTelBindLogConsume(c chan *model.UserTel) {
	filterStart := 1536572121
	filterEnd := 1536616436
	for {
		a, ok := <-c
		if !ok {
			log.Error("asoTelBindLogChan closed")
			return
		}
		for {
			var (
				err         error
				telBindTime int64
			)
			if telBindTime, err = s.d.AsoTelBindLog(context.Background(), a.Mid); err != nil {
				log.Error("fail to get AsoTelBindLog error(%+v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if telBindTime == 0 {
				break
			}
			if telBindTime > int64(filterStart) && telBindTime < int64(filterEnd) {
				break
			}
			userTel := &model.UserTel{
				Mid:         a.Mid,
				TelBindTime: telBindTime,
			}
			if _, err = s.d.UpdateUserTelBindTime(context.Background(), userTel); err != nil {
				log.Error("fail to update tel bind log userTel(%+v) error(%+v)", userTel, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}
