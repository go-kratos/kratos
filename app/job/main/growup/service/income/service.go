package income

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/dao/dataplatform"
	"go-common/app/job/main/growup/dao/email"
	incomeD "go-common/app/job/main/growup/dao/income"
	"go-common/app/job/main/growup/service"
	"go-common/app/job/main/growup/service/ctrl"
	"go-common/library/log"
)

// Service struct
type Service struct {
	conf     *conf.Config
	dao      *incomeD.Dao
	avCharge *AvChargeSvr
	income   *Income
	ratio    *ChargeRatioSvr
	email    *email.Dao
	dp       *dataplatform.Dao
}

// New fn
func New(c *conf.Config, executor ctrl.Executor) (s *Service) {
	s = &Service{
		conf:  c,
		dao:   incomeD.New(c),
		email: email.New(c),
		dp:    dataplatform.New(c),
	}
	s.avCharge = NewAvChargeSvr(s.dao)
	s.income = NewIncome(batchSize, s.dao)
	s.ratio = NewChargeRatioSvr(s.dao)
	log.Info("income service start")

	executor.Submit(
		s.calDailyCreativeIncome,
		s.syncBGM,
	)
	return s
}

func (s *Service) syncBGM(ctx context.Context) {
	for {
		time.Sleep(service.NextDay(1, 0, 0))
		msg := ""
		log.Info("sync BGM begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		err := s.SyncBgmInfo(context.TODO())
		if err != nil {
			log.Error("s.SyncBgmInfo error(%v)", err)
		}
		if err != nil {
			msg = fmt.Sprintf("SyncBgmInfo error(%v)", err)
		} else {
			msg = "Success"
		}
		err = s.email.SendMail(time.Now(), msg, "创作激励同步bgm%d年%d月%d日", "shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com")
		if err != nil {
			log.Error("s.email.SendMail error(%v)", err)
		}
		log.Info("sync BGM end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (s *Service) calDailyCreativeIncome(ctx context.Context) {
	for {
		time.Sleep(service.NextDay(18, 0, 0))
		log.Info("calDailyCreativeIncome begin:%v", time.Now().Format("2006-01-02 15:04:05"))
		now := time.Now()
		date := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
		err := s.RunAndSendMail(context.TODO(), date)
		if err != nil {
			log.Error("s.RunAndSendMail error(%v)", err)
		}
		if err == nil {
			err = s.RunStatis(context.TODO(), date)
			if err != nil {
				log.Error("s.RunStatis error(%v)", err)
			}
		}
		log.Info("calDailyCreativeIncome end:%v", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}
