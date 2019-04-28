package service

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
)

//LogCredit log credit log
func (s *Service) LogCredit(c context.Context, arg *upcrmmodel.ArgCreditLogAdd) (err error) {
	err = s.upcrmdb.AddLog(arg)
	if err != nil {
		log.Error("fail to log credit log, log=%+v, err=%v", arg, err)
		return
	}
	log.Info("log credit log to db, {%+v}", arg)
	return
}

//WriteStatData write log to db
func (s *Service) WriteStatData() {
	defer func() {
		s.wg.Done()
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("write stat data Runtime error caught, try recover: %+v", r)
			s.wg.Add(1)
			go s.WriteStatData()
		}
	}()
	for s.running {
		select {
		case creditScore := <-s.CreditScoreInputChan:
			var err = s.upcrmdb.AddOrUpdateCreditScore(creditScore)
			if err != nil {
				log.Error("fail to insert credit score, mid=%d", creditScore.Mid)
				continue
			}
			var affectedRow int64
			affectedRow, err = s.upcrmdb.UpdateCreditScore(creditScore.Score, creditScore.Mid)
			log.Info("insert credit score, mid=%d, update base info, affected=%d, err=%+v", creditScore.Mid, affectedRow, err)
		case <-s.closeChan:
			log.Info("server closing, close routine")
		}
	}
}
